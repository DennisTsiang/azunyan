package manager

import (
	"github.com/callummance/azunyan/db"
	"github.com/callummance/azunyan/models"
	broadcast "github.com/dustin/go-broadcast"
)

const broadcastBufLen int = 10

type BroadcastData struct {
	Name    string
	Content interface{}
}

func InitBroadcaster() broadcast.Broadcaster {
	return broadcast.NewBroadcaster(broadcastBufLen)
}

func (m *KaraokeManager) SubscribeToChanges() chan interface{} {
	newListener := make(chan interface{}, broadcastBufLen)
	m.UpdateSubscribers.Register(newListener)

	return newListener
}

func (m *KaraokeManager) Unsubscribe(target chan interface{}) {
	m.UpdateSubscribers.Unregister(target)
}

func (m *KaraokeManager) SendBroadcast(bc BroadcastData) {
	m.UpdateSubscribers.Submit(bc)
}

/*FetchAndUpdateListenersQueue fetches the complete and partialqueues and
broadcasts them to the listeners.
mode specifies the mode that the listener should use:
0 - Default, reuse the partial divs
1 - Clears the divs in the partialQueue and recreates new divs
*/
func FetchAndUpdateListenersQueue(m *KaraokeManager, mode int) error {
	completeQueue, partialQueue, err := GetQueue(m)
	if err != nil {
		m.Logger.Printf("Failed to get song queue state due to error %v", err)
		return err
	}
	UpdateListenersQueue(m, completeQueue, partialQueue, mode)
	return nil
}

func UpdateListenersQueue(m *KaraokeManager, complete []models.QueueItem, partial []models.QueueItem, mode int) {
	//Send updates
	m.SendBroadcast(BroadcastData{
		Name: "queue",
		Content: struct {
			Mode   int
			Queues map[string][]models.QueueItem
		}{
			mode,
			map[string][]models.QueueItem{
				"complete": complete,
				"partial":  partial,
			},
		},
	})
}

func UpdateListenersCur(m *KaraokeManager, cur *models.QueueItem) {
	//Send updates
	m.SendBroadcast(BroadcastData{
		Name:    "cur",
		Content: cur,
	})
}

func FetchAndUpdateListenersCur(m *KaraokeManager) error {
	state, err := db.GetEngineState(m, m.Config.KaraokeConfig.SessionName)
	if err != nil {
		m.Logger.Printf("Failed to get session data due to error %q", err)
		return err
	}
	cur := state.NowPlaying
	UpdateListenersCur(m, cur)
	return nil
}

// Broadcasts an update to the "singers" listeners
func UpdateNumberSingers(m *KaraokeManager, singers int) {
	m.SendBroadcast(BroadcastData{
		Name:    "singers",
		Content: singers,
	})
}
