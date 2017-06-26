package events

import (
    "testing"
    "encoding/hex"
    "encoding/base64"
)

func TestEventID_Parse(t *testing.T) {
    in := "952822DE6A627EA459E1E7A8964191C79FCCFB14EA545D93741B5CF3ED71A09A"
    id := NewEventID()
    id.Parse(in)

    expected := [32]byte{
        149, 40, 34, 222, 106, 98, 126, 164,
        89, 225, 231, 168, 150, 65, 145, 199,
        159, 204, 251, 20, 234, 84, 93, 147,
        116, 27, 92, 243, 237, 113, 160, 154 }

    if [32]byte(id) != expected {
        t.Errorf("Parsed ID did not match expected. Actual: %x; Expected: %x", [32]byte(id), expected)
    }
}

func TestEvent_ID(t *testing.T) {
    prevId := NewEventID()
    err := prevId.Parse("952822DE6A627EA459E1E7A8964191C79FCCFB14EA545D93741B5CF3ED71A09A")
    if err != nil {
        t.Fatal("Error parsing test ID.", err.Error())
    }

    data, err := base64.StdEncoding.DecodeString("eyJhY2NvdW50SWQiOiI0Y2FmYmU2Yy1kMzYxLTRiZTMtYjcyZS1kNjNhNDQzNmUyMDQiLCJhbW91bnQiOjEwMDAwfQ==")
    if err != nil {
        t.Fatal("Error parsing test data.", err.Error())
    }

    e := &Event {
        PreviousEvent: prevId,
            Type: "AmountDeposited",
        Data: data,
    }

    id := e.ID()

    expected := "1fbf23ba470895ccdba7a401b51a939539052af72b301b86558e845880756493"
    result := hex.EncodeToString(id[:])

    if result != expected {
        t.Errorf("Generated ID did not match expected. Result: %s; Expected: %s", result, expected)
    }
}
