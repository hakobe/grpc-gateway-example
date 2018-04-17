package main

import (
	"context"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	pb_articles "github.com/hakobe/grpc-gateway-example/articles"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type article struct {
	title   string
	body    string
	created time.Time
}

type server struct {
	articles []*article
	mutex    *sync.RWMutex
}

func (s *server) Post(ctx context.Context, postRequest *pb_articles.PostRequest) (*pb_articles.PostResponse, error) {
	a := postRequest.GetArticle()
	if a == nil {
		return nil, errors.New("No article")
	}
	created, err := ptypes.Timestamp(a.Created)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot build timestamp from article")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.articles = append(s.articles, &article{
		title:   a.Title,
		body:    a.Body,
		created: created,
	})
	sort.Slice(s.articles, func(i, j int) bool { return s.articles[i].created.Before(s.articles[j].created) })

	return &pb_articles.PostResponse{Article: a}, nil
}

func (s *server) Recent(ctx context.Context, _empty *empty.Empty) (*pb_articles.RecentResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	results := make([]*pb_articles.Article, 0)
	for i := 0; i < 100 && i < len(s.articles); i++ {
		a := s.articles[i]
		created, err := ptypes.TimestampProto(a.created)
		if err != nil {
			return nil, errors.Wrap(err, "Cannot build time from article")
		}
		results = append(results, &pb_articles.Article{
			Title:   a.title,
			Body:    a.body,
			Created: created,
		})
	}

	return &pb_articles.RecentResponse{Articles: results}, nil
}

func serve(s *server, hostPort string) {
	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb_articles.RegisterArticlesServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	hostPort := os.Getenv("HOST_PORT")
	if hostPort == "" {
		hostPort = "0.0.0.0:5000"
	}

	s := &server{
		articles: make([]*article, 0),
		mutex:    &sync.RWMutex{},
	}

	log.Println("Starting a server on " + hostPort)
	serve(s, hostPort)
}
