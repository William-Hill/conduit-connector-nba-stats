package nbastats

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type NBAStatsQueryParams struct {
	College          string // Default is empty
	Conference       string // Default is empty
	Country          string // Default is empty
	DateFrom         string // Default is empty
	DateTo           string // Default is empty
	Division         string // Default is empty
	DraftPick        string // Default is empty
	DraftYear        string // Default is empty
	GameScope        string // Default is empty
	Height           string // Default is empty
	ISTRound         string // Default is empty
	LastNGames       int    // Default is 0
	LeagueID         string // Default is "00"
	Location         string // Default is empty
	Month            int    // Default is 0
	OpponentTeamID   int    // Default is 0
	Outcome          string // Default is empty
	PORound          int    // Default is 0
	PerMode          string // Default is "PerGame"
	PlayerExperience string // Default is empty
	PlayerOrTeam     string // Default is "Player"
	PlayerPosition   string // Default is empty
	PtMeasureType    string // Default is "SpeedDistance"
	Season           string // Default is "2023-24"
	SeasonSegment    string // Default is empty
	SeasonType       string // Default is "Regular Season"
	StarterBench     string // Default is empty
	TeamID           int    // Default is 0
	VsConference     string // Default is empty
	VsDivision       string // Default is empty
	Weight           string // Default is empty
}

// NewNBAStatsQueryParams creates an instance of NBAStatsQueryParams with default values.
func NewNBAStatsQueryParams() NBAStatsQueryParams {
	return NBAStatsQueryParams{
		LastNGames:     0,
		LeagueID:       "00",
		Month:          0,
		OpponentTeamID: 0,
		PORound:        0,
		PerMode:        "PerGame",
		PlayerOrTeam:   "Player",
		PtMeasureType:  "SpeedDistance",
		Season:         "2023-24",
		SeasonType:     "Regular Season",
		TeamID:         0,
	}
}

func buildNBAStatsURL(params NBAStatsQueryParams) string {
	baseURL := "https://stats.nba.com/stats/leaguedashptstats"
	values := url.Values{}

	values.Set("College", params.College)
	values.Set("Conference", params.Conference)
	values.Set("Country", params.Country)
	values.Set("DateFrom", params.DateFrom)
	values.Set("DateTo", params.DateTo)
	values.Set("Division", params.Division)
	values.Set("DraftPick", params.DraftPick)
	values.Set("DraftYear", params.DraftYear)
	values.Set("GameScope", params.GameScope)
	values.Set("Height", params.Height)
	values.Set("ISTRound", params.ISTRound)
	values.Set("LastNGames", strconv.Itoa(params.LastNGames))
	values.Set("LeagueID", params.LeagueID)
	values.Set("Location", params.Location)
	values.Set("Month", strconv.Itoa(params.Month))
	values.Set("OpponentTeamID", strconv.Itoa(params.OpponentTeamID))
	values.Set("Outcome", params.Outcome)
	values.Set("PORound", strconv.Itoa(params.PORound))
	values.Set("PerMode", params.PerMode)
	values.Set("PlayerExperience", params.PlayerExperience)
	values.Set("PlayerOrTeam", params.PlayerOrTeam)
	values.Set("PlayerPosition", params.PlayerPosition)
	values.Set("PtMeasureType", params.PtMeasureType)
	values.Set("Season", params.Season)
	values.Set("SeasonSegment", params.SeasonSegment)
	values.Set("SeasonType", params.SeasonType)
	values.Set("StarterBench", params.StarterBench)
	values.Set("TeamID", strconv.Itoa(params.TeamID))
	values.Set("VsConference", params.VsConference)
	values.Set("VsDivision", params.VsDivision)
	values.Set("Weight", params.Weight)

	return baseURL + "?" + values.Encode()
}

// ResponseData structure reflects the JSON structure of the API response.
type ResponseData struct {
	Resource   string      `json:"resource"`
	Parameters interface{} `json:"parameters"`
	ResultSets []struct {
		Name    string          `json:"name"`
		Headers []string        `json:"headers"`
		RowSet  [][]interface{} `json:"rowSet"`
	} `json:"resultSets"`
}

func fetchNBASpeedDistanceStats(perMode string) ([]byte, error) {
	// url := "https://stats.nba.com/stats/leaguedashptstats?College=&Conference=&Country=&DateFrom=&DateTo=&Division=&DraftPick=&DraftYear=&GameScope=&Height=&ISTRound=&LastNGames=0&LeagueID=00&Location=&Month=0&OpponentTeamID=0&Outcome=&PORound=0&PerMode=PerGame&PlayerExperience=&PlayerOrTeam=Player&PlayerPosition=&PtMeasureType=SpeedDistance&Season=2023-24&SeasonSegment=&SeasonType=Regular%20Season&StarterBench=&TeamID=0&VsConference=&VsDivision=&Weight="

	var nbaStatsQuery = NewNBAStatsQueryParams()
	nbaStatsQuery.PerMode = perMode
	url := buildNBAStatsURL(nbaStatsQuery)
	fmt.Printf("url: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Set the required headers
	req.Header.Set("Host", "stats.nba.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:72.0) Gecko/20100101 Firefox/72.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// Note: The "Accept-Encoding" header is managed by the http.Client. If you set it manually, you must also handle the encoding yourself.
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("x-nba-stats-origin", "stats")
	req.Header.Set("x-nba-stats-token", "true")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://stats.nba.com/")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	// Create a new client with a timeout
	client := &http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for status code 200 OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error: status code %d\n", resp.StatusCode)
		return nil, err
	}

	// Handle gzip encoding
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// Read and print response body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, err
}
