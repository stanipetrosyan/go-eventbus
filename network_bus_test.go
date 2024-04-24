package goeventbus

/*
func TestNetworkBus(t *testing.T) {
	var wg sync.WaitGroup

	t.Run("should send a message to the srever using client", func(t *testing.T) {
		wg.Add(1)

		eventbus := NewEventBus()

		eventbus.Channel("my-channel").Subscriber().Listen(func(context Context) {
			assert.Equal(t, "Hello there", context.Result().Data)
			wg.Done()
		})

		var networkBus = NewNetworkBus(NewEventBus(), "localhost:8080", "/")

		go networkBus.Server().Listen()

		go networkBus.Client().Connect()

		wg.Wait()
	})

}
*/
