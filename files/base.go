package files

import "go.uber.org/zap"

type usecase struct {
	fileRepo    FileRepository
	fileStorage FileStorageAdapter
	txManager   TransactionManager
	logger      *zap.Logger
	cfg         Config
}
