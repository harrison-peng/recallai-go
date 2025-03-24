package recallaigo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type BotService interface {
	ListBots(ctx context.Context, params *ListBotsParams) (*ListBotResponse, error)
	CreateBot(ctx context.Context, request *CreateBotRequest) (*Bot, error)
	ListChatMessages(ctx context.Context, botID string, params ...ListChatMessagesParams) (*ListMessagesResponse, error)
	RetrieveBot(ctx context.Context, botID string) (*Bot, error)
	UpdateScheduledBot(ctx context.Context, botID string, request *Bot) (*Bot, error)
	DeleteScheduledBot(ctx context.Context, botID string) error
	DeleteBotMedia(ctx context.Context, botID string) error
	GetBotLogs(ctx context.Context, botID string) (*LogEntry, error)
	OutputAudio(ctx context.Context, botID string, request *OutputAudioRequest) (*Bot, error)
	StopOutputAudio(ctx context.Context, botID string) error
	OutputMedia(ctx context.Context, botID string, request *OutputMedia) (*Bot, error)
	StopOutputMedia(ctx context.Context, botID string) error
	StartScreenshare(ctx context.Context, botID string, request *OutputVideoRequest) (*Bot, error)
	StopScreenshare(ctx context.Context, botID string) error
	OutputVideo(ctx context.Context, botID string, request *OutputVideoRequest) (*Bot, error)
	StopOutputVideo(ctx context.Context, botID string) error
	PauseRecording(ctx context.Context, botID string) (*Bot, error)
	RequestRecordingPermission(ctx context.Context, botID string) (*Bot, error)
	ResumeRecording(ctx context.Context, botID string) (*Bot, error)
	SendChatMessage(ctx context.Context, botID string, request *SendChatMessageRequest) (*Bot, error)
	GetSpeakerTimeline(ctx context.Context, botID string, params ...GetSpeakerTimelineParams) ([]SpeakerTimelineEntry, error)
	StartRecording(ctx context.Context, botID string, request *StartRecordingRequest) (*Bot, error)
	StopRecording(ctx context.Context, botID string) (*Bot, error)
	GetBotTranscript(ctx context.Context, botID string, params ...GetBotTranscriptParams) ([]TranscriptEntry, error)
	AnalyzeBotMedia(ctx context.Context, botId string, request *AnalyzeBotMediaRequest) (*AnalyzeBotMediaResponse, error)
}

type BotClient struct {
	client *Client
}

type Platform string

const (
	PlatformZoom                Platform = "zoom"
	PlatformGoogleMeet          Platform = "google_meet"
	PlatformGotoMeeting         Platform = "goto_meeting"
	PlatformMicrosoftTeams      Platform = "microsoft_teams"
	PlatformMicrosoftTeamsLive  Platform = "microsoft_teams_live"
	PlatformWebex               Platform = "webex"
	PlatformChimeSdk            Platform = "chime_sdk"
	PlatformSlackAuthenticator  Platform = "slack_authenticator"
	PlatformSlackHuddleObserver Platform = "slack_huddle_observer"
)

func (p Platform) String() string {
	return string(p)
}

type Status string

const (
	StatusReady                      Status = "ready"
	StatusJoiningCall                Status = "joining_call"
	StatusInWaitingRoom              Status = "in_waiting_room"
	StatusInCallNotRecording         Status = "in_call_not_recording"
	StatusRecordingPermissionAllowed Status = "recording_permission_allowed"
	StatusRecordingPermissionDenied  Status = "recording_permission_denied"
	StatusInCallRecording            Status = "in_call_recording"
	StatusRecordingDone              Status = "recording_done"
	StatusCallEnded                  Status = "call_ended"
	StatusDone                       Status = "done"
	StatusFatal                      Status = "fatal"
	StatusMediaExpired               Status = "media_expired"
	StatusAnalysisDone               Status = "analysis_done"
	StatusAnalysisFailed             Status = "analysis_failed"
)

func (s Status) String() string {
	return string(s)
}

// ListBotsParams defines the parameters for filtering and paginating the list of bots.
type ListBotsParams struct {
	// Filter bots that joined after this date-time (ISO 8601 format)
	JoinAtAfter string `json:"join_at_after,omitempty"`
	// Filter bots that joined before this date-time (ISO 8601 format)
	JoinAtBefore string `json:"join_at_before,omitempty"`
	// Filter bots by the meeting URL
	MeetingURL string `json:"meeting_url,omitempty"`
	// Specify the page number for pagination
	Page int `json:"page,omitempty"`
	// Filter bots by platform(s)
	Platform []Platform `json:"platform,omitempty"`
	// Filter bots by status(es)
	Status []Status `json:"status,omitempty"`
}

// ListBotResponse represents the response body for the List method
type ListBotResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Bot  `json:"results"`
}

// Bot represents a bot with its configuration and state.
type Bot struct {
	ID string `json:"id"`
	// The url of the meeting.
	MeetingURL MeetingURL `json:"meeting_url"`
	// The name of the bot that will be displayed in the call.
	// (Note: Authenticated Google Meet bots will use the Google account name and this field will be ignored.)
	// The string length should be ≤ 100. Defaults to "Meeting Notetaker".
	BotName string `json:"bot_name" validate:"max=100,default=Meeting Notetaker"`
	// The time at which the bot will join the call, formatted in ISO 8601.
	// This field can only be read from scheduled bots that have not yet joined a call.
	// Once a bot has joined a call, its join_at will be cleared.
	JoinAt              *string              `json:"join_at,omitempty"`
	VideoURL            string               `json:"video_url"`
	MediaRetentionEnd   string               `json:"media_retention_end"`
	StatusChanges       []StatusChange       `json:"status_changes"`
	MeetingMetadata     MeetingMetadata      `json:"meeting_metadata"`
	MeetingParticipants []MeetingParticipant `json:"meeting_participants"`
	// The settings for real-time transcription.
	RealTimeTranscription *RealTimeTranscription `json:"real_time_transcription,omitempty"`
	// The settings for real-time media output.
	RealTimeMedia *RealTimeMedia `json:"real_time_media,omitempty"`
	// The options for transcription settings.
	TranscriptionOptions *TranscriptionOptions `json:"transcription_options,omitempty"`
	// The mode in which the recording will be made. Defaults to "speaker_view".
	RecordingMode RecordingMode `json:"recording_mode" validate:"oneof=speaker_view gallery_view gallery_view_v2 audio_only,default=speaker_view"`
	// Additional options for recording mode.
	RecordingModeOptions *RecordingModeOptions `json:"recording_mode_options,omitempty"`
	// Settings to include the bot in the recording.
	IncludeBotInRecording *IncludeBotInRecording `json:"include_bot_in_recording,omitempty"`
	Recordings            []Recording            `json:"recordings"`
	// Settings for the bot output media.
	OutputMedia *OutputMedia `json:"output_media,omitempty"`
	// Settings for the bot to output video. Image should be 16:9. Recommended resolution is 640x360.
	AutomaticVideoOutput *AutomaticVideoOutput `json:"automatic_video_output,omitempty"`
	// (BETA) Settings for the bot to output audio.
	AutomaticAudioOutput *AutomaticAudioOutput `json:"automatic_audio_output,omitempty"`
	// (BETA) Settings for the bot to send chat messages.
	// (Note: Chat functionality is only supported for Zoom, Google Meet, and Microsoft Teams currently.)
	Chat *Chat `json:"chat,omitempty"`
	// (BETA) Settings for the bot to automatically leave the meeting.
	AutomaticLeave *AutomaticLeave `json:"automatic_leave,omitempty"`
	// Configure bot variants per meeting platforms, e.g. {"zoom": "web_4_core"}.
	Variant          *Variant          `json:"variants,omitempty"`
	CalendarMeetings []CalendarMeeting `json:"calendar_meetings"`
	// Zoom specific parameters
	Zoom *Zoom `json:"zoom,omitempty"`
	// Google Meet specific parameters
	GoogleMeet *GoogleMeet `json:"google_meet,omitempty"`
	// Slack Authenticator specific parameters
	SlackAuthenticator *SlackAuthenticator `json:"slack_authenticator,omitempty"`
	// Slack Huddle Observer specific parameters
	SlackHuddleObserver *SlackHuddleObserver `json:"slack_huddle_observer,omitempty"`
	// Metadata for the bot, which can include additional information as key-value pairs.
	Metadata  map[string]string `json:"metadata,omitempty"`
	Recording string            `json:"recording"`
}

type MeetingURL struct {
	MeetingID       string  `json:"meeting_id"`
	MeetingPassword string  `json:"meeting_password"`
	TK              *string `json:"tk"`
	Platform        string  `json:"platform"`
}

type RecordingMode string

const (
	SpeakerView   RecordingMode = "speaker_view"
	GalleryView   RecordingMode = "gallery_view"
	GalleryViewV2 RecordingMode = "gallery_view_v2"
	AudioOnly     RecordingMode = "audio_only"
)

type MeetingMetadata struct {
	Title           string `json:"title"`
	ZoomMeetingUUID string `json:"zoom_meeting_uuid"`
	SlackChannelID  string `json:"slack_channel_id"`
	SlackHuddleID   string `json:"slack_huddle_id"`
}

type MeetingParticipant struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Events []struct {
		Code      string `json:"code"`
		CreatedAt string `json:"created_at"`
	} `json:"events"`
	IsHost    bool   `json:"is_host"`
	Platform  string `json:"platform"`
	ExtraData struct {
		Zoom struct {
			UserGUID   string `json:"user_guid"`
			Guest      bool   `json:"guest"`
			ConfUserID string `json:"conf_user_id"`
		} `json:"zoom"`
		MicrosoftTeams struct {
			ParticipantType string `json:"participant_type"`
			Role            string `json:"role"`
			MeetingRole     string `json:"meeting_role"`
			UserID          string `json:"user_id"`
			TenantID        string `json:"tenant_id"`
			ClientVersion   string `json:"client_version"`
		} `json:"microsoft_teams"`
		Slack struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
		} `json:"slack"`
	} `json:"extra_data"`
}

type StatusChange struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	SubCode   string `json:"sub_code"`
}

type TranscriptionProvider string

const (
	TranscriptionProviderDeepgram               TranscriptionProvider = "deepgram"
	TranscriptionProviderAssemblyAIAsyncChunked TranscriptionProvider = "assembly_ai_async_chunked"
	TranscriptionProviderAssemblyAI             TranscriptionProvider = "assembly_ai"
	TranscriptionProviderRev                    TranscriptionProvider = "rev"
	TranscriptionProviderAWSTranscribe          TranscriptionProvider = "aws_transcribe"
	TranscriptionProviderSpeechmatics           TranscriptionProvider = "speechmatics"
	TranscriptionProviderGladia                 TranscriptionProvider = "gladia"
	TranscriptionProviderGladiaV2               TranscriptionProvider = "gladia_v2"
	TranscriptionProviderMeetingCaptions        TranscriptionProvider = "meeting_captions"
	TranscriptionProviderNone                   TranscriptionProvider = "none"
)

type TranscriptionOptions struct {
	Provider TranscriptionProvider `json:"provider"`
}

type IncludeBotInRecording struct {
	Audio bool `json:"audio"`
}

type RealTimeTranscription struct {
	DestinationURL      string `json:"destination_url"`
	PartialResults      bool   `json:"partial_results"`
	EnhancedDiarization bool   `json:"enhanced_diarization"`
}

type RealTimeMedia struct {
	RTMPDestinationURL                         string `json:"rtmp_destination_url"`
	WebsocketVideoDestinationURL               string `json:"websocket_video_destination_url"`
	WebsocketAudioDestinationURL               string `json:"websocket_audio_destination_url"`
	WebsocketSpeakerTimelineDestinationURL     string `json:"websocket_speaker_timeline_destination_url"`
	WebsocketSpeakerTimelineExcludeNullSpeaker bool   `json:"websocket_speaker_timeline_exclude_null_speaker"`
	WebhookCallEventsDestinationURL            string `json:"webhook_call_events_destination_url"`
	WebhookChatMessagesDestinationURL          string `json:"webhook_chat_messages_destination_url"`
}

type RecordingModeOptions struct {
	ParticipantVideoWhenScreenshare string `json:"participant_video_when_screenshare"`
	StartRecordingOn                string `json:"start_recording_on"`
}

type Chat struct {
	OnBotJoin         ChatOnBotJoin         `json:"on_bot_join"`
	OnParticipantJoin ChatOnParticipantJoin `json:"on_participant_join"`
}

type ChatOnBotJoin struct {
	SendTo  string `json:"send_to"`
	Message string `json:"message"`
	Pin     bool   `json:"pin"`
}

type ChatOnParticipantJoin struct {
	Message     string `json:"message"`
	ExcludeHost bool   `json:"exclude_host"`
}

type Recording struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type OutputMedia struct {
	Camera      OutputMediaSetting `json:"camera"`
	Screenshare OutputMediaSetting `json:"screenshare"`
}

type OutputMediaSetting struct {
	Kind   string            `json:"kind"`
	Config OutputMediaConfig `json:"config"`
}

type OutputMediaConfig struct {
	URL string `json:"url"`
}

type AutomaticVideoOutput struct {
	InCallRecording    AutomaticVideoOutputConfig `json:"in_call_recording"`
	InCallNotRecording AutomaticVideoOutputConfig `json:"in_call_not_recording"`
}

type AutomaticVideoOutputConfig struct {
	Kind string `json:"kind"`
}

type AutomaticAudioOutput struct {
	InCallRecording InCallRecording `json:"in_call_recording"`
}

type InCallRecording struct {
	Data                    InCallRecordingData     `json:"data"`
	ReplayOnParticipantJoin ReplayOnParticipantJoin `json:"replay_on_participant_join"`
}

type InCallRecordingData struct {
	Kind string `json:"kind"`
}

type ReplayOnParticipantJoin struct {
	DebounceMode     string `json:"debounce_mode"`
	DebounceInterval int    `json:"debounce_interval"`
	DisableAfter     int    `json:"disable_after"`
}

type AutomaticLeave struct {
	WaitingRoomTimeout               int              `json:"waiting_room_timeout"`
	NooneJoinedTimeout               int              `json:"noone_joined_timeout"`
	EveryoneLeftTimeout              int              `json:"everyone_left_timeout"`
	InCallNotRecordingTimeout        int              `json:"in_call_not_recording_timeout"`
	InCallRecordingTimeout           int              `json:"in_call_recording_timeout"`
	RecordingPermissionDeniedTimeout int              `json:"recording_permission_denied_timeout"`
	SilenceDetection                 SilenceDetection `json:"silence_detection"`
	BotDetection                     BotDetection     `json:"bot_detection"`
}

type SilenceDetection struct {
	Timeout       int `json:"timeout"`
	ActivateAfter int `json:"activate_after"`
}

type BotDetection struct {
	UsingParticipantEvents UsingParticipantEvents `json:"using_participant_events"`
	UsingParticipantNames  UsingParticipantNames  `json:"using_participant_names"`
}

type UsingParticipantEvents struct {
	Timeout       int `json:"timeout"`
	ActivateAfter int `json:"activate_after"`
}

type UsingParticipantNames struct {
	Timeout       int      `json:"timeout"`
	ActivateAfter int      `json:"activate_after"`
	Matches       []string `json:"matches"`
}

type Variant struct {
	Zoom           VariantOption `json:"zoom"`
	GoogleMeet     VariantOption `json:"google_meet"`
	MicrosoftTeams VariantOption `json:"microsoft_teams"`
}

type VariantOption string

const (
	VariantWeb      VariantOption = "web"
	VariantWeb4Core VariantOption = "web_4_core"
	VariantNative   VariantOption = "native"
)

type CalendarMeeting struct {
	ID           string       `json:"id"`
	StartTime    string       `json:"start_time"`
	EndTime      string       `json:"end_time"`
	CalendarUser CalendarUser `json:"calendar_user"`
}

type CalendarUser struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
}

type Zoom struct {
	JoinTokenURL string `json:"join_token_url"`
	ZakURL       string `json:"zak_url"`
	UserEmail    string `json:"user_email"`
}

type GoogleMeet struct {
	LoginRequired      bool   `json:"login_required"`
	GoogleLoginGroupID string `json:"google_login_group_id"`
}

type SlackAuthenticator struct {
	SlackTeamIntegrationID string `json:"slack_team_integration_id"`
	TeamDomain             string `json:"team_domain"`
	LoginEmail             string `json:"login_email"`
	ProfilePhotoBase64JPG  string `json:"profile_photo_base64_jpg"`
}

type SlackHuddleObserver struct {
	SlackTeamIntegrationID    string   `json:"slack_team_integration_id"`
	TeamDomain                string   `json:"team_domain"`
	LoginEmail                string   `json:"login_email"`
	AutoJoinPublicHuddles     bool     `json:"auto_join_public_huddles"`
	AskToJoinPrivateHuddles   bool     `json:"ask_to_join_private_huddles"`
	AskToJoinMessage          string   `json:"ask_to_join_message"`
	FilterHuddlesByUserEmails []string `json:"filter_huddles_by_user_emails"`
	ProfilePhotoBase64JPG     string   `json:"profile_photo_base64_jpg"`
	HuddleBotAPIToken         string   `json:"huddle_bot_api_token"`
}

type Message struct {
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	To        string `json:"to"`
	Sender    Sender `json:"sender"`
}

type Sender struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsHost    bool      `json:"is_host"`
	Platform  string    `json:"platform"`
	ExtraData ExtraData `json:"extra_data"`
}

type ExtraData struct {
	Zoom           ZoomData           `json:"zoom"`
	MicrosoftTeams MicrosoftTeamsData `json:"microsoft_teams"`
	Slack          SlackData          `json:"slack"`
}

type ZoomData struct {
	UserGUID   string `json:"user_guid"`
	Guest      bool   `json:"guest"`
	ConfUserID string `json:"conf_user_id"`
}

type MicrosoftTeamsData struct {
	ParticipantType string `json:"participant_type"`
	Role            string `json:"role"`
	MeetingRole     string `json:"meeting_role"`
	UserID          string `json:"user_id"`
	TenantID        string `json:"tenant_id"`
	ClientVersion   string `json:"client_version"`
}

type SlackData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// Get a list of all bots
// see https://docs.recall.ai/reference/bot_list
func (c *BotClient) ListBots(ctx context.Context, params *ListBotsParams) (*ListBotResponse, error) {
	queryParams := buildQueryParams(params)

	res, err := c.client.request(ctx, http.MethodGet, "bot", queryParams, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to list bots: %w", err)
	}
	defer res.Body.Close()

	// bodyBytes, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read response body: %w", err)
	// }
	// fmt.Println(string(bodyBytes))

	var response ListBotResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func buildQueryParams(params *ListBotsParams) map[string][]string {
	queryParams := make(map[string][]string)

	if params == nil {
		return queryParams
	}

	addQueryParam := func(key, value string) {
		if value != "" {
			queryParams[key] = []string{value}
		}
	}

	addQueryParam("join_at_after", params.JoinAtAfter)
	addQueryParam("join_at_before", params.JoinAtBefore)
	addQueryParam("meeting_url", params.MeetingURL)

	if params.Page != 0 {
		queryParams["page"] = []string{fmt.Sprintf("%d", params.Page)}
	}
	if len(params.Platform) > 0 {
		queryParams["platform"] = convertToStringSlice(params.Platform)
	}
	if len(params.Status) > 0 {
		queryParams["status"] = convertToStringSlice(params.Status)
	}

	return queryParams
}

func convertToStringSlice[T fmt.Stringer](items []T) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.String()
	}
	return result
}

// CreateBotRequest represents the request body for the CreateBot method
type CreateBotRequest struct {
	// The url of the meeting. For example, https://zoom.us/j/123?pwd=456. This field will be cleared a few days after the bot has joined a call.
	MeetingURL string `json:"meeting_url"`
	// The name of the bot that will be displayed in the call.
	// (Note: Authenticated Google Meet bots will use the Google account name and this field will be ignored.)
	// The string length should be ≤ 100. Defaults to "Meeting Notetaker".
	BotName string `json:"bot_name" validate:"max=100,default=Meeting Notetaker"`
	// The time at which the bot will join the call, formatted in ISO 8601.
	// This field can only be read from scheduled bots that have not yet joined a call.
	// Once a bot has joined a call, its join_at will be cleared.
	JoinAt *string `json:"join_at,omitempty"`
	// The settings for real-time transcription.
	RealTimeTranscription *RealTimeTranscription `json:"real_time_transcription,omitempty"`
	// The settings for real-time media output.
	RealTimeMedia *RealTimeMedia `json:"real_time_media,omitempty"`
	// The options for transcription settings.
	TranscriptionOptions *TranscriptionOptions `json:"transcription_options,omitempty"`
	// The mode in which the recording will be made. Defaults to "speaker_view".
	RecordingMode RecordingMode `json:"recording_mode" validate:"oneof=speaker_view gallery_view gallery_view_v2 audio_only,default=speaker_view"`
	// Additional options for recording mode.
	RecordingModeOptions *RecordingModeOptions `json:"recording_mode_options,omitempty"`
	// Settings to include the bot in the recording.
	IncludeBotInRecording *IncludeBotInRecording `json:"include_bot_in_recording,omitempty"`
	Recordings            []Recording            `json:"recordings"`
	// Settings for the bot output media.
	OutputMedia *OutputMedia `json:"output_media,omitempty"`
	// Settings for the bot to output video. Image should be 16:9. Recommended resolution is 640x360.
	AutomaticVideoOutput *AutomaticVideoOutput `json:"automatic_video_output,omitempty"`
	// (BETA) Settings for the bot to output audio.
	AutomaticAudioOutput *AutomaticAudioOutput `json:"automatic_audio_output,omitempty"`
	// (BETA) Settings for the bot to send chat messages.
	// (Note: Chat functionality is only supported for Zoom, Google Meet, and Microsoft Teams currently.)
	Chat *Chat `json:"chat,omitempty"`
	// (BETA) Settings for the bot to automatically leave the meeting.
	AutomaticLeave *AutomaticLeave `json:"automatic_leave,omitempty"`
	// Configure bot variants per meeting platforms, e.g. {"zoom": "web_4_core"}.
	Variant *Variant `json:"variants,omitempty"`
	// Zoom specific parameters
	Zoom *Zoom `json:"zoom,omitempty"`
	// Google Meet specific parameters
	GoogleMeet *GoogleMeet `json:"google_meet,omitempty"`
	// Slack Authenticator specific parameters
	SlackAuthenticator *SlackAuthenticator `json:"slack_authenticator,omitempty"`
	// Slack Huddle Observer specific parameters
	SlackHuddleObserver *SlackHuddleObserver `json:"slack_huddle_observer,omitempty"`
	// Metadata for the bot, which can include additional information as key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (r *CreateBotRequest) Validate() error {
	if r.MeetingURL == "" {
		return fmt.Errorf("meeting URL is required")
	}
	if r.BotName == "" {
		return fmt.Errorf("bot name is required")
	}

	return nil
}

// CreateBot a new bot
// see https://docs.recall.ai/reference/bot_create
func (c *BotClient) CreateBot(ctx context.Context, request *CreateBotRequest) (*Bot, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	res, err := c.client.request(ctx, http.MethodPost, "bot", nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}
	defer res.Body.Close()

	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

type ListChatMessagesParams struct {
	Cursor   string
	Ordering string
}

type ListMessagesResponse struct {
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Message `json:"results"`
}

// Get list of chat messages read by the bot in the meeting(excluding messages sent by the bot itself).
// see https://docs.recall.ai/reference/bot_chat_messages_list
func (c *BotClient) ListChatMessages(ctx context.Context, botID string, params ...ListChatMessagesParams) (*ListMessagesResponse, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/chat-messages", botID)

	// Prepare query parameters
	queryParams := make(map[string][]string)
	if len(params) > 0 {
		param := params[0]
		if param.Cursor != "" {
			queryParams["cursor"] = []string{param.Cursor}
		}
		if param.Ordering != "" {
			queryParams["ordering"] = []string{param.Ordering}
		}
	}

	// Make the request
	res, err := c.client.request(ctx, http.MethodGet, path, queryParams, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat messages: %w", err)
	}
	defer res.Body.Close()

	// Decode the response
	var message ListMessagesResponse
	if err := json.NewDecoder(res.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &message, nil
}

// RetrieveBot retrieves a bot by its ID.
// see https://docs.recall.ai/reference/bot_retrieve
func (c *BotClient) RetrieveBot(ctx context.Context, botID string) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s", botID)

	// Make the request
	res, err := c.client.request(ctx, http.MethodGet, path, nil, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bot: %w", err)
	}
	defer res.Body.Close()

	// Decode the response
	var bot Bot
	if err := json.NewDecoder(res.Body).Decode(&bot); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &bot, nil
}

// UpdateScheduledBot updates the schedule of a bot by its ID.
// see https://docs.recall.ai/reference/bot_partial_update
func (c *BotClient) UpdateScheduledBot(ctx context.Context, botID string, request *Bot) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s", botID)

	// Make the request
	res, err := c.client.request(ctx, http.MethodPatch, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to update scheduled bot: %w", err)
	}
	defer res.Body.Close()

	// Decode the response
	var bot Bot
	if err := json.NewDecoder(res.Body).Decode(&bot); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &bot, nil
}

// DeleteScheduledBot deletes a bot by its ID.
// see https://docs.recall.ai/reference/bot_destroy
func (c *BotClient) DeleteScheduledBot(ctx context.Context, botID string) error {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s", botID)

	// Make the request
	res, err := c.client.request(ctx, http.MethodDelete, path, nil, nil, apiVersionV1)
	if err != nil {
		return fmt.Errorf("failed to delete scheduled bot: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// DeleteBotMedia deletes the media of a bot by its ID.
// see https://docs.recall.ai/reference/bot_delete_media_create
func (c *BotClient) DeleteBotMedia(ctx context.Context, botID string) error {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/delete_media", botID)

	// Make the request
	res, err := c.client.request(ctx, http.MethodPost, path, nil, nil, apiVersionV1)
	if err != nil {
		return fmt.Errorf("failed to delete bot media: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// Get the results of additional analysis specified by the intelligence parameter. If the call is not yet complete, this returns results from any real-time analysis performed so-far.
// Not implemented yet
// see https://docs.recall.ai/reference/bot_intelligence_retrieve
// func (c *BotClient) GetBotIntelligence(ctx context.Context, botID string) (*IntelligenceResult, error) {
// 	// TODO: Implement this method
// 	return nil, nil
// }

// RemoveBotFromCall removes the bot from a call by its ID.
// This action is irreversible.
// see https://docs.recall.ai/reference/bot_leave_call_create
// Not implemented yet
// func (c *BotClient) RemoveBotFromCall(ctx context.Context, botID string) error {
// 	// TODO: Implement this method
// 	return nil
// }

// LogEntry represents a single log entry with level, message, and created_at fields.
type LogEntry struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// GetBotLogs retrieves the logs produced by the bot by its ID.
// see https://docs.recall.ai/reference/bot_logs_retrieve
func (c *BotClient) GetBotLogs(ctx context.Context, botID string) (*LogEntry, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/logs", botID)

	// Make the request
	res, err := c.client.request(ctx, http.MethodGet, path, nil, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot logs: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a slice of LogEntry
	var log LogEntry
	if err := json.NewDecoder(res.Body).Decode(&log); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &log, nil
}

type OutputAudioKind string

const (
	OutputAudioKindMp3 OutputAudioKind = "mp3"
)

// OutputAudioRequest represents the request body for the OutputAudio method.
type OutputAudioRequest struct {
	Kind    OutputAudioKind `json:"kind" `
	B64Data string          `json:"b64_data"`
}

// OutputAudio causes the bot to output audio.
// see https://docs.recall.ai/reference/bot_output_audio_create
func (c *BotClient) OutputAudio(ctx context.Context, botID string, request *OutputAudioRequest) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_audio", botID)

	// Make the request with the provided OutputAudioRequest
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to output audio: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// StopOutputAudio stops the bot from outputting audio.
// see https://docs.recall.ai/reference/bot_output_audio_destroy
func (c *BotClient) StopOutputAudio(ctx context.Context, botID string) error {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_audio", botID)

	// Make the DELETE request to stop outputting audio
	res, err := c.client.request(ctx, http.MethodDelete, path, nil, nil, apiVersionV1)
	if err != nil {
		return fmt.Errorf("failed to stop output audio: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// OutputMedia causes the bot to start outputting media.
// see https://docs.recall.ai/reference/bot_output_media_create
func (c *BotClient) OutputMedia(ctx context.Context, botID string, request *OutputMedia) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_media", botID)

	// Make the request with the provided OutputMediaRequest
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to output media: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// StopOutputMedia stops the bot from outputting media.
// see https://docs.recall.ai/reference/bot_output_media_destroy
func (c *BotClient) StopOutputMedia(ctx context.Context, botID string) error {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_media", botID)

	// Make the DELETE request to stop outputting media
	res, err := c.client.request(ctx, http.MethodDelete, path, nil, nil, apiVersionV1)
	if err != nil {
		return fmt.Errorf("failed to stop output media: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

type OutputVideoKind string

const (
	OutputVideoKindJpeg OutputAudioKind = "jpeg"
)

// OutputAudioRequest represents the request body for the OutputAudio method.
type OutputVideoRequest struct {
	Kind    OutputVideoKind `json:"kind" `
	B64Data string          `json:"b64_data"`
}

// StartScreenshare causes the bot to start screensharing.
// see https://docs.recall.ai/reference/bot_output_screenshare_create
func (c *BotClient) StartScreenshare(ctx context.Context, botID string, request *OutputVideoRequest) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_screenshare", botID)

	// Make the POST request with the provided OutputVideoRequest
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to start screenshare: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// StopScreenshare causes the bot to stop screensharing.
// see https://docs.recall.ai/reference/bot_output_screenshare_destroy
func (c *BotClient) StopScreenshare(ctx context.Context, botID string) error {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_screenshare", botID)

	// Make the DELETE request to stop screensharing
	res, err := c.client.request(ctx, http.MethodDelete, path, nil, nil, apiVersionV1)
	if err != nil {
		return fmt.Errorf("failed to stop screenshare: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// OutputVideo causes the bot to start outputting video.
// see https://docs.recall.ai/reference/bot_output_video_create
func (c *BotClient) OutputVideo(ctx context.Context, botID string, request *OutputVideoRequest) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_video", botID)

	// Make the POST request with the provided OutputVideoRequest
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to output video: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// StopOutputVideo stops the bot from outputting video.
// see https://docs.recall.ai/reference/bot_output_video_destroy
func (c *BotClient) StopOutputVideo(ctx context.Context, botID string) error {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/output_video", botID)

	// Make the DELETE request to stop outputting video
	res, err := c.client.request(ctx, http.MethodDelete, path, nil, nil, apiVersionV1)
	if err != nil {
		return fmt.Errorf("failed to stop output video: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// PauseRecording instructs the bot to pause the current recording.
// see https://docs.recall.ai/reference/bot_pause_recording_create
func (c *BotClient) PauseRecording(ctx context.Context, botID string) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/pause_recording", botID)

	// Make the POST request to pause the recording
	res, err := c.client.request(ctx, http.MethodPost, path, nil, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to pause recording: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// RequestRecordingPermission requests recording permission from the host.
// This is applicable for Zoom only.
// see https://docs.recall.ai/reference/bot_request_recording_permission_create
func (c *BotClient) RequestRecordingPermission(ctx context.Context, botID string) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/request_recording_permission", botID)

	// Make the POST request to request recording permission
	res, err := c.client.request(ctx, http.MethodPost, path, nil, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to request recording permission: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// ResumeRecording resumes a paused recording for the bot.
// see https://docs.recall.ai/reference/bot_resume_recording_create
func (c *BotClient) ResumeRecording(ctx context.Context, botID string) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/resume_recording", botID)

	// Make the POST request to resume the recording
	res, err := c.client.request(ctx, http.MethodPost, path, nil, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to resume recording: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

type SendChatMessageRequest struct {
	// The person or group that the message will be sent to. On non-Zoom platforms, "everyone" is currently the only supported option, meaning the message will be sent to everyone in the meeting.
	To string `json:"to,omitempty"`

	// The message content. Required, with a maximum length of 4096 characters.
	Message string `json:"message"`

	// Whether to pin the message. Defaults to false.
	Pin bool `json:"pin,omitempty"`
}

// SendChatMessage sends a message in the meeting chat.
// see https://docs.recall.ai/reference/bot_send_chat_message_create
func (c *BotClient) SendChatMessage(ctx context.Context, botID string, request *SendChatMessageRequest) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/send_chat_message", botID)

	// Make the POST request to send the chat message
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to send chat message: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

type GetSpeakerTimelineParams struct {
	ExcludeNullSpeaker bool
}

// SpeakerTimelineEntry represents a single entry in the speaker timeline.
type SpeakerTimelineEntry struct {
	Name      string  `json:"name"`
	UserID    int     `json:"user_id"`
	Timestamp float64 `json:"timestamp"`
}

// GetSpeakerTimeline retrieves the speaker timeline produced by the bot.
// If the call is not yet complete, this returns the speaker timeline so-far.
// see https://docs.recall.ai/reference/bot_speaker_timeline_list
func (c *BotClient) GetSpeakerTimeline(ctx context.Context, botID string, params ...GetSpeakerTimelineParams) ([]SpeakerTimelineEntry, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/speaker_timeline", botID)

	// Prepare query parameters
	queryParams := make(map[string][]string)
	if len(params) > 0 && params[0].ExcludeNullSpeaker {
		queryParams["exclude_null_speaker"] = []string{"true"}
	}

	// Make the GET request to retrieve the speaker timeline
	res, err := c.client.request(ctx, http.MethodGet, path, queryParams, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to get speaker timeline: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a slice of SpeakerTimelineEntry
	var timeline []SpeakerTimelineEntry
	if err := json.NewDecoder(res.Body).Decode(&timeline); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return timeline, nil
}

// StartRecordingRequest represents the request body for the StartRecording method.
type StartRecordingRequest struct {
	RecordingMode         RecordingMode         `json:"recording_mode"`
	RecordingModeOptions  RecordingModeOptions  `json:"recording_mode_options"`
	RealTimeTranscription RealTimeTranscription `json:"real_time_transcription"`
	RealTimeMedia         RealTimeMedia         `json:"real_time_media"`
	TranscriptionOptions  TranscriptionOptions  `json:"transcription_options"`
}

// StartRecording instructs the bot to start recording the meeting.
// This will restart the current recording if one is already in progress.
// see https://docs.recall.ai/reference/bot_start_recording_create
func (c *BotClient) StartRecording(ctx context.Context, botID string, request *StartRecordingRequest) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/start_recording", botID)

	// Make the POST request with the provided StartRecordingRequest
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to start recording: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// StopRecording instructs the bot to stop recording the meeting.
// see https://docs.recall.ai/reference/bot_stop_recording_create
func (c *BotClient) StopRecording(ctx context.Context, botID string) (*Bot, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/stop_recording", botID)

	// Make the POST request to stop recording
	res, err := c.client.request(ctx, http.MethodPost, path, nil, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to stop recording: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a Bot
	var response Bot
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetTranscriptParams represents the query parameters for the GetTranscript method.
type GetBotTranscriptParams struct {
	EnhancedDiarization bool
}

// TranscriptEntry represents a single entry in the bot's transcript.
type TranscriptEntry struct {
	Speaker   string       `json:"speaker"`
	SpeakerID int          `json:"speaker_id"`
	Language  string       `json:"language"`
	Words     []WordDetail `json:"words"`
}

// WordDetail represents the details of a word in the transcript.
type WordDetail struct {
	Text           string  `json:"text"`
	StartTimestamp float64 `json:"start_timestamp"`
	EndTimestamp   float64 `json:"end_timestamp"`
	Language       string  `json:"language"`
	Confidence     float64 `json:"confidence"`
}

// GetBotTranscript retrieves the transcript produced by the bot by its ID.
// see https://docs.recall.ai/reference/bot_transcript_list
func (c *BotClient) GetBotTranscript(ctx context.Context, botID string, params ...GetBotTranscriptParams) ([]TranscriptEntry, error) {
	// Construct the URL path with the bot_id
	path := fmt.Sprintf("bot/%s/transcript", botID)

	// Prepare query parameters
	queryParams := make(map[string][]string)
	if len(params) > 0 && params[0].EnhancedDiarization {
		queryParams["enhanced_diarization"] = []string{"true"}
	}

	// Make the GET request with the query parameters
	res, err := c.client.request(ctx, http.MethodGet, path, queryParams, nil, apiVersionV1)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot transcript: %w", err)
	}
	defer res.Body.Close()

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response body into a slice of TranscriptEntry
	var transcript []TranscriptEntry
	if err := json.NewDecoder(res.Body).Decode(&transcript); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return transcript, nil
}

type AnalyzeBotMediaRequest struct {
	// Transcription requests for various services.
	AssemblyAIAsyncTranscription   AssemblyAIAsyncTranscription   `json:"assemblyai_async_transcription"`
	SpeechmaticsAsyncTranscription SpeechmaticsAsyncTranscription `json:"speechmatics_async_transcription"`
	RevAsyncTranscription          RevAsyncTranscription          `json:"rev_async_transcription"`
	DeepgramAsyncTranscription     DeepgramAsyncTranscription     `json:"deepgram_async_transcription"`
}

// AssemblyAIAsyncTranscription represents the request for asynchronous transcription using AssemblyAI.
type AssemblyAIAsyncTranscription struct {
	// Language specifies the language of the audio.
	// Refer to: https://www.assemblyai.com/docs/speech-to-text/supported-languages
	Language string `json:"language"`

	// AudioEndAt indicates the timestamp to stop processing the audio.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.audio_end_at
	AudioEndAt int `json:"audio_end_at"`

	// AudioStartFrom indicates the timestamp to start processing the audio.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.audio_start_from
	AudioStartFrom int `json:"audio_start_from"`

	// AutoChapters enables automatic chapter detection in the transcription.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.auto_chapters
	AutoChapters bool `json:"auto_chapters"`

	// AutoHighlights enables automatic highlights detection in the transcription.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.auto_highlights
	AutoHighlights bool `json:"auto_highlights"`

	// BoostParam specifies the boost parameter to enhance transcription accuracy.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.boost_param
	BoostParam string `json:"boost_param"`

	// ContentSafety enables analysis for potentially sensitive content.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.content_safety
	ContentSafety bool `json:"content_safety"`

	// ContentSafetyConfidence sets the confidence threshold for content safety analysis.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.content_safety_confidence
	ContentSafetyConfidence int `json:"content_safety_confidence"`

	// CustomSpelling allows configuration of custom spelling for specific words.
	// Refer to: https://www.assemblyai.com/docs/api-reference/transcripts/submit#request.body.custom_spelling
	CustomSpelling []map[string]string `json:"custom_spelling"`
}

// SpeechmaticsAsyncTranscription represents the request for asynchronous transcription using Speechmatics.
type SpeechmaticsAsyncTranscription struct {
	// Language specifies the language model to process the audio input.
	// This field is required and normally specified as an ISO language code.
	Language string `json:"language"`

	// Domain requests a specialized model based on 'language' but optimized for a particular field.
	// Examples include "finance" or "medical".
	Domain string `json:"domain"`

	// OutputLocale specifies the language locale to be used when generating the transcription output.
	// Normally specified as an ISO language code.
	OutputLocale string `json:"output_locale"`

	// OperatingPoint specifies the operating point for the transcription.
	// Enum: "standardenhanced"
	OperatingPoint string `json:"operating_point"`

	// AdditionalVocab is a list of custom words or phrases that should be recognized.
	// Alternative pronunciations can be specified to aid recognition.
	AdditionalVocab []map[string]string `json:"additional_vocab"`
}

// RevAsyncTranscription represents the request for asynchronous transcription using Rev.
type RevAsyncTranscription struct {
	// DetectLanguage enables automatic language detection.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=detect_language&t=request
	DetectLanguage bool `json:"detect_language"`

	// Language specifies the language model to process the audio input.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=language&t=request
	Language string `json:"language"`

	// SkipDiarization skips the speaker diarization process.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=skip_diarization&t=request
	SkipDiarization bool `json:"skip_diarization"`

	// SkipPostprocessing skips the post-processing step.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=skip_postprocessing&t=request
	SkipPostprocessing bool `json:"skip_postprocessing"`

	// SkipPunctuation skips the punctuation step in the transcription.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=skip_punctuation&t=request
	SkipPunctuation bool `json:"skip_punctuation"`

	// RemoveDisfluencies removes disfluencies from the transcription.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=remove_disfluencies&t=request
	RemoveDisfluencies bool `json:"remove_disfluencies"`

	// RemoveAtmospherics removes atmospheric sounds from the transcription.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=remove_atmospherics&t=request
	RemoveAtmospherics bool `json:"remove_atmospherics"`

	// FilterProfanity filters out profane language from the transcription.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=filter_profanity&t=request
	FilterProfanity bool `json:"filter_profanity"`

	// CustomVocabularyID specifies the ID of a custom vocabulary to use.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=custom_vocabulary_id&t=request
	CustomVocabularyID string `json:"custom_vocabulary_id"`

	// CustomVocabularies is a list of custom vocabulary strings to be used.
	// Docs: https://docs.rev.ai/api/asynchronous/reference/#operation/SubmitTranscriptionJob!ct=application/json&path=custom_vocabularies&t=request
	CustomVocabularies []string `json:"custom_vocabularies"`
}

// DeepgramAsyncTranscription represents the request for asynchronous transcription using Deepgram.
type DeepgramAsyncTranscription struct {
	// Tier specifies the tier of the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/tier/
	Tier string `json:"tier"`

	// Model specifies the model to use for the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/model/
	Model string `json:"model"`

	// Version specifies the version of the model to use.
	// Docs: https://developers.deepgram.com/documentation/features/version/
	Version string `json:"version"`

	// Language specifies the language model to process the audio input.
	// Docs: https://developers.deepgram.com/documentation/features/language/
	Language string `json:"language"`

	// DetectLanguage enables automatic language detection.
	// Docs: https://developers.deepgram.com/documentation/features/detect-language/
	DetectLanguage bool `json:"detect_language"`

	// Punctuate enables automatic punctuation in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/punctuate/
	Punctuate bool `json:"punctuate"`

	// ProfanityFilter filters out profane language from the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/profanity-filter/
	ProfanityFilter bool `json:"profanity_filter"`

	// Redact specifies the list of strings to redact from the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/redact/
	Redact []string `json:"redact"`

	// Diarize enables speaker diarization in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/diarize/
	Diarize bool `json:"diarize"`

	// DiarizeVersion specifies the version of the diarization model to use.
	// Docs: https://developers.deepgram.com/documentation/features/diarize/
	DiarizeVersion string `json:"diarize_version"`

	// SmartFormat enables smart formatting in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/smart-format/
	SmartFormat bool `json:"smart_format"`

	// Numerals enables numeral formatting in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/numerals/
	Numerals bool `json:"numerals"`

	// Search specifies the list of strings to search for in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/search/
	Search []string `json:"search"`

	// Replace specifies the list of strings to replace in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/replace/
	Replace []string `json:"replace"`

	// Keywords specifies the list of keywords to highlight in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/keywords/
	Keywords []string `json:"keywords"`

	// Summarize enables summarization of the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/summarize/
	Summarize bool `json:"summarize"`

	// DetectTopics enables topic detection in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/detect-topics/
	DetectTopics bool `json:"detect_topics"`

	// Tag enables tagging in the transcription.
	// Docs: https://developers.deepgram.com/documentation/features/tag/
	Tag bool `json:"tag"`

	// CredentialID specifies the ID of the Deepgram credential to use for this transcription.
	// If not specified, the default credential will be used.
	CredentialID string `json:"credential_id"`
}

// GladiaV2AsyncTranscription represents the request for asynchronous transcription using Gladia V2.
type GladiaV2AsyncTranscription struct {
	// ContextPrompt is a string used for context in transcription.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#context-prompt
	ContextPrompt string `json:"context_prompt"`

	// CustomVocabulary can be a boolean or a list of vocabulary items.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#custom-vocabulary
	CustomVocabulary interface{} `json:"custom_vocabulary"` // Can be bool or list

	// CustomVocabularyConfig holds configuration for custom vocabulary.
	CustomVocabularyConfig CustomVocabularyConfig `json:"custom_vocabulary_config"`

	// EnableCodeSwitching allows detection of multiple languages.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#multiple-languages-detection-code-switching
	EnableCodeSwitching bool `json:"enable_code_switching"`

	// CodeSwitchingConfig holds configuration for guided code switching.
	CodeSwitchingConfig CodeSwitchingConfig `json:"code_switching_config"`

	// Subtitles enables the export of caption files.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#export-srt-or-vtt-caption-files
	Subtitles bool `json:"subtitles"`

	// SubtitlesConfig holds configuration for subtitles.
	SubtitlesConfig SubtitlesConfig `json:"subtitles_config"`

	// DiarizationConfig holds configuration for speaker diarization.
	DiarizationConfig DiarizationConfig `json:"diarization_config"`

	// TranslationConfig holds configuration for translation.
	TranslationConfig TranslationConfig `json:"translation_config"`

	// SummarizationConfig holds configuration for summarization.
	SummarizationConfig SummarizationConfig `json:"summarization_config"`

	// Moderation enables content moderation.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/moderation
	Moderation bool `json:"moderation"`

	// NamedEntityRecognition enables named entity recognition.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/named%20entity%20recognition
	NamedEntityRecognition bool `json:"named_entity_recognition"`

	// Chapterization enables chapterization of the transcription.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/chapterization
	Chapterization bool `json:"chapterization"`

	// NameConsistency ensures consistent naming in transcription.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#name-consistency
	NameConsistency bool `json:"name_consistency"`

	// CustomSpelling enables custom spelling in transcription.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#custom-spelling
	CustomSpelling bool `json:"custom_spelling"`

	// CustomSpellingConfig holds configuration for custom spelling.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#custom-spelling
	CustomSpellingConfig string `json:"custom_spelling_config"`

	// StructuredDataExtraction enables extraction of structured data.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/structured%20data%20extraction
	StructuredDataExtraction bool `json:"structured_data_extraction"`

	// StructuredDataExtractionConfig holds configuration for structured data extraction.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/structured%20data%20extraction
	StructuredDataExtractionConfig string `json:"structured_data_extraction_config"`

	// SentimentAnalysis enables sentiment analysis.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/sentiment%20analysis
	SentimentAnalysis bool `json:"sentiment_analysis"`

	// AudioToLLM enables audio to language model processing.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/audio%20to%20llm
	AudioToLLM bool `json:"audio_to_llm"`

	// AudioToLLMConfig holds configuration for audio to language model processing.
	AudioToLLMConfig AudioToLLMConfig `json:"audio_to_llm_config"`

	// CustomMetadata allows adding custom metadata to transcription.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#adding-custom-metadata
	CustomMetadata string `json:"custom_metadata"`

	// Sentences enables sentence-level transcription.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#sentences
	Sentences bool `json:"sentences"`

	// DisplayMode enables display mode in transcription.
	// Docs: https://docs.gladia.io/api-reference/api-v2/Transcription/post-v2transcription
	DisplayMode bool `json:"display_mode"`

	// PunctuationEnhanced enables enhanced punctuation in transcription.
	// Add comment: This field is used to enhance punctuation in the transcription.
	PunctuationEnhanced bool `json:"punctuation_enhanced"`
}

// CustomVocabularyConfig holds configuration for custom vocabulary.
type CustomVocabularyConfig struct {
	// DefaultIntensity sets the default intensity for custom vocabulary.
	// This field is a double.
	DefaultIntensity float64 `json:"default_intensity"`

	// Vocabulary is a list of custom vocabulary.
	// This field is an array of objects.
	Vocabulary []Vocabulary `json:"vocabulary"`
}

// Vocabulary represents a custom vocabulary item.
type Vocabulary struct {
	// Value to add to custom vocabulary.
	// This field is required.
	Value string `json:"value"`

	// Intensity for this specific custom vocabulary item (0 to 1).
	// This field is a double.
	Intensity float64 `json:"intensity"`

	// Pronunciations are optional pronunciations for this vocabulary item.
	// This field is an array of strings.
	Pronunciations []string `json:"pronunciations,omitempty"`

	// Language is an optional language for this vocabulary item.
	// Defaults to transcription language if not specified.
	Language string `json:"language,omitempty"`
}

// CodeSwitchingConfig holds configuration for guided code switching.
type CodeSwitchingConfig struct {
	// Languages is a list of languages supported for code switching.
	// This field is required and is an array of strings.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/languages
	Languages []string `json:"languages"`
}

// SubtitlesConfig holds configuration for subtitles.
type SubtitlesConfig struct {
	// Formats specifies the subtitle file formats to export.
	// This field is required and is an array of strings.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speech-recognition#export-srt-or-vtt-caption-files
	Formats []string `json:"formats"`
}

// DiarizationConfig holds configuration for speaker diarization.
type DiarizationConfig struct {
	// NumberOfSpeakers specifies the expected number of speakers in the audio.
	// This field is an integer.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speaker-diarization#improving-diarization-accuracy
	NumberOfSpeakers int `json:"number_of_speakers"`

	// MinSpeakers specifies the minimum number of speakers to be detected.
	// This field is an integer.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speaker-diarization#improving-diarization-accuracy
	MinSpeakers int `json:"min_speakers"`

	// MaxSpeakers specifies the maximum number of speakers to be detected.
	// This field is an integer.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/speaker-diarization#improving-diarization-accuracy
	MaxSpeakers int `json:"max_speakers"`
}

// TranslationConfig holds configuration for translation.
type TranslationConfig struct {
	// TargetLanguages specifies the languages to which the audio should be translated.
	// This field is required and is an array of strings.
	// Docs: https://docs.gladia.io/chapters/speech-to-text-api/pages/languages
	TargetLanguages []string `json:"target_languages"`

	// Model specifies the translation model to be used.
	// This field is a string.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/translation
	Model string `json:"model"`
}

// SummarizationConfig holds configuration for summarization.
type SummarizationConfig struct {
	// Type specifies the type of summarization to be performed.
	// This field is a string.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/summarization
	Type string `json:"type"`
}

// AudioToLLMConfig holds configuration for audio to language model processing.
type AudioToLLMConfig struct {
	// Prompts specifies the prompts to be used for audio to language model processing.
	// This field is required and is an array of strings.
	// Docs: https://docs.gladia.io/chapters/audio-intelligence/pages/audio%20to%20llm
	Prompts []string `json:"prompts"`
}

// AnalyzeBotMediaResponse represents the response for the AnalyzeBotMedia method.
type AnalyzeBotMediaResponse struct {
	JobId string `json:"job_id"`
}

// AnalyzeBotMedia runs analysis on the bot's media.
// Not implemented yet
// see https://docs.recall.ai/reference/bot_analyze_create
func (c *BotClient) AnalyzeBotMedia(ctx context.Context, botId string, request *AnalyzeBotMediaRequest) (*AnalyzeBotMediaResponse, error) {
	path := fmt.Sprintf("bot/%s/analyze", botId)

	// Make the POST request to analyze bot media
	res, err := c.client.request(ctx, http.MethodPost, path, nil, request, apiVersionV2Beta)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze bot media: %w", err)
	}
	defer res.Body.Close()

	// Decode the response into AnalyzeBotMediaResponse
	var response AnalyzeBotMediaResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
