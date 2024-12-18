package service

import (
	"context"
	"errors"
	"testing"

	"awesomeProject/internal/entites"
	"awesomeProject/internal/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateTokens(t *testing.T) {
	type mockBehavior func(s *mocks.MockAuth, userID string)

	testCases := []struct {
		name            string
		userID          string
		clientIP        string
		mockBehavior    mockBehavior
		expectedJWT     string
		expectedRefresh string
		expectedError   error
	}{
		{
			name:     "Success",
			userID:   "user123",
			clientIP: "127.0.0.1",
			mockBehavior: func(s *mocks.MockAuth, userID string) {
				// Ожидание вызова GetUserById
				s.EXPECT().GetUserById(gomock.Any(), userID).Return(&entites.User{Email: "test@example.com"}, nil).AnyTimes()
				// Ожидание вызова InsertUserInfo
				s.EXPECT().InsertUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedJWT:     "mockedJWT",
			expectedRefresh: "mockedRefresh",
		},
		{
			name:     "User not found",
			userID:   "user123",
			clientIP: "127.0.0.1",
			mockBehavior: func(s *mocks.MockAuth, userID string) {
				// Ожидание вызова GetUserById с отсутствующим пользователем
				s.EXPECT().GetUserById(gomock.Any(), userID).Return(nil, nil).AnyTimes()
			},
			expectedError: errors.New("no refresh token found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := mocks.NewMockAuth(c)
			authService := NewAuthService(store, "mockSecret", "mockSecret")

			tc.mockBehavior(store, tc.userID)

			jwt, refresh, err := authService.GenerateTokens(context.Background(), tc.userID, tc.clientIP)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, jwt)
				assert.NotEmpty(t, refresh)
			}
		})
	}
}
