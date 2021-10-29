test_integration:
	 go test -tags=integration -v -cover ./...

test_unit:
	go test -v -cover ./...

mock_user:
	 mockgen -source services/user/service.go -destination services/user/mock_user_service.go -package user

.PHONY: test_integration test_unit mock_user