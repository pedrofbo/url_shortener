package internal

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pedrofbo/url_shortener/dynamodb"
)

type Item struct {
	ShortUrl string `json:"short_url" dynamodbav:"short_url"`
	LongUrl  string `json:"long_url" dynamodbav:"long_url"`
}

type ParticleEntry struct {
	Particle string `json:"particle" dynamodbav:"particle"`
}

func GetLongUrl(tableName string, shortUrl string) (*Item, error) {
	input := map[string]types.AttributeValue{
		"short_url": &types.AttributeValueMemberS{Value: shortUrl},
	}
	result, err := dynamodb.GetItem(tableName, input)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch item: %s", err)
	}
	item := &Item{}
	err = attributevalue.UnmarshalMap(result, item)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Record, %v", err)
	}

	return item, nil
}

func GenerateShortUrl(longUrl string, particleTableName string) (*string, error) {
	particles, err := getAllParticles(particleTableName)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch particles from table %s", particleTableName)
	}
	if particles == nil {
		return nil, fmt.Errorf("No particles found on table %s", particleTableName)
	}

	value_1, value_2, value_3 := generateSampleValues(longUrl, len(*particles))
	shortUrl := ((*particles)[value_1].Particle + (*particles)[value_2].Particle + (*particles)[value_3].Particle)
	return &shortUrl, nil
}

func getAllUrls(tableName string) (*[]Item, error) {
	result, err := dynamodb.GetAllItems(tableName)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch all items from table %s: %s", tableName, err)
	}
	items := &[]Item{}
	err = attributevalue.UnmarshalListOfMaps(result, items)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Record, %v", err)
	}

	return items, nil
}

func getAllParticles(tableName string) (*[]ParticleEntry, error) {
	result, err := dynamodb.GetAllItems(tableName)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch all particles from table %s: %s", tableName, err)
	}
	items := &[]ParticleEntry{}
	err = attributevalue.UnmarshalListOfMaps(result, items)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal record, %v", err)
	}

	return items, nil
}

func generateSampleValues(input string, maxValue int) (int, int, int) {
	seed := getSeedFromString(input)
	rand.Seed(seed)
	return rand.Intn(maxValue), rand.Intn(maxValue), rand.Intn(maxValue)
}

func getSeedFromString(input string) int64 {
	h := sha256.New()
	io.WriteString(h, input)
	return int64(binary.BigEndian.Uint64(h.Sum(nil)))
}

func JoinUrl(baseUrl string, endpoint string) (string, error) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		return "", nil
	}
	url.Path = path.Join(url.Path, endpoint)
	return url.String(), nil
}
