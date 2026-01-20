/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, filenames ...string) *mongo.Client {
	if err := godotenv.Load(filenames...); err != nil {
		log.Error(err)
	}

	uri := os.Getenv("mongodb_uri")
	if uri == "" {
		log.Fatal("You must set your 'mongodb_uri' environment variable.")
	}
	clientOptions := options.Client().ApplyURI(uri).SetTimeout(time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	return client
}

func Disconnect(client *mongo.Client, ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
