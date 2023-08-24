package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/olezhek28/microservices_course/week_8/chat/pkg/chat_v1"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := desc.NewChatV1Client(conn)

	// Создаем новый чат на сервере
	chatID, err := createChat(ctx, client, []string{"oleg", "ivan"})
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	log.Printf(fmt.Sprintf("%s: %s\n", color.GreenString("Chat created"), color.YellowString(chatID)))

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Подключаемся к чату от имени пользователя oleg
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "oleg")
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	// Подключаемся к чату от имени пользователя ivan
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "ivan")
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	wg.Wait()
}

func connectChat(ctx context.Context, client desc.ChatV1Client, chatID string, username string) error {
	stream, err := client.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format(time.RFC3339)),
				color.BlueString(message.GetFrom()),
				message.GetText(),
			)
		}
	}()

	for {
		// Ниже пример того, как можно считывать сообщения из консоли
		// в демонстрационных целях будем засылать в чат рандомный текст раз в 5 секунд
		//scanner := bufio.NewScanner(os.Stdin)
		//var lines strings.Builder
		//
		//for {
		//	scanner.Scan()
		//	line := scanner.Text()
		//	if len(line) == 0 {
		//		break
		//	}
		//
		//	lines.WriteString(line)
		//	lines.WriteString("\n")
		//}
		//
		//err = scanner.Err()
		//if err != nil {
		//	log.Println("failed to scan message: ", err)
		//}

		time.Sleep(5 * time.Second)

		text := gofakeit.Word()

		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			ChatId: chatID,
			Message: &desc.Message{
				From:      username,
				Text:      text,
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func createChat(ctx context.Context, client desc.ChatV1Client, usernames []string) (string, error) {
	res, err := client.CreateChat(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return res.GetChatId(), nil
}
