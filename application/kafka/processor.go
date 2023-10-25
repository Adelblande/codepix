package kafka

import (
	"fmt"
	"os"

	"github.com/adelblande/codepix/application/dto"
	"github.com/adelblande/codepix/application/factory"
	"github.com/adelblande/codepix/application/usecase"
	"github.com/adelblande/codepix/domain/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
)

type KafkaProcessor struct {
	Database *gorm.DB
	Producer *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database: database,
		Producer: producer,
		DeliveryChan: deliveryChan,
	}
}

func (k *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id": os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}

	c, err := ckafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}
	
	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionTopicConfirmation")}
	c.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionTopics := "transactions"
	transactionTopicsConfirmation := "transactions_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionTopics:
		k.processTransaction(msg)
	case transactionTopicsConfirmation:
		k.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transactionDto := dto.NewTransactionDto()
	err := transactionDto.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	createdTransaction, err := transactionUseCase.Register(
		transactionDto.AccountId,
		transactionDto.Amount,
		transactionDto.PixKeyTo,
		transactionDto.PixKeyKindTo,
		transactionDto.Description,
	)

	if err != nil {
		return err
	}

	topic := "bank"+createdTransaction.PixKeyTo.Account.Bank.Code
	transactionDto.ID = createdTransaction.ID
	transactionDto.Status = model.TransactionPending

	transactionJson, err := transactionDto.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transactionDto := dto.NewTransactionDto()
	err := transactionDto.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	if transactionDto.Status == model.TransactionConfirmed {

		err = k.confirmTransaction(transactionDto, transactionUseCase)
		if err != nil {
			return err
		}
	}

	if transactionDto.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transactionDto.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(
	transactionDto *dto.TransactionDto, 
	transactionUseCase usecase.TransactionUseCase,
) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transactionDto.ID)
	if err != nil {
		return err
	}

	topic := "bank"+confirmedTransaction.AccountFrom.Bank.ID
	transactionJson, err := transactionDto.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}
