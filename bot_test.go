package recallaigo_test

import (
	"context"
	"net/http"
	"testing"

	recallaigo "github.com/harrison-peng/recallai-go"
)

func TestBotClient(t *testing.T) {
	t.Run("ListBots", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			len        int
			wantErr    bool
		}{
			{
				name:       "returns bots",
				statusCode: http.StatusOK,
				filePath:   "test_data/list_bots.json",
				len:        1,
			},
			{
				name:       "returns error",
				statusCode: http.StatusBadRequest,
				filePath:   "test_data/error.json",
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := recallaigo.NewClient("some_token", recallaigo.WithHTTPClient(c))
				got, err := client.Bot.ListBots(context.Background(), nil)

				if (err != nil) != tt.wantErr {
					t.Errorf("ListBots() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.len != 0 && len(got.Results) != tt.len {
					t.Errorf("ListBots got %d, want: %d", len(got.Results), tt.len)
				}

				if tt.wantErr && err == nil {
					t.Error("ListBots() error = nil, wantErr true")
				}
			})
		}
	})

	t.Run("CreateBot", func(t *testing.T) {
		tests := []struct {
			name       string
			statusCode int
			filePath   string
			wantErr    bool
		}{
			{
				name:       "returns bot",
				statusCode: http.StatusOK,
				filePath:   "test_data/create_bot.json",
			},
			{
				name:       "returns error",
				statusCode: http.StatusBadRequest,
				filePath:   "test_data/error.json",
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)

				request := recallaigo.CreateBotRequest{
					MeetingURL:    "https://test.com",
					BotName:       "Test Bot",
					RecordingMode: recallaigo.SpeakerView,
				}

				client := recallaigo.NewClient("some_token", recallaigo.WithHTTPClient(c))
				res, err := client.Bot.CreateBot(context.Background(), &request)

				if tt.wantErr {
					if err == nil {
						t.Error("CreateBot() error = nil, wantErr true")
					}
				} else {
					if err != nil {
						t.Errorf("CreateBot() error = %v, wantErr %v", err, tt.wantErr)
					}

					if res == nil {
						t.Error("CreateBot() response = nil, wantErr false")
					}
				}
			})
		}
	})

}
