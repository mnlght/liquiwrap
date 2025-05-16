# liquiwrap

A set of tools to retrieve tournament and match information from liquipedia.net. Written for https://gorn.pw/.

```http
    go get github.com/mnlght/liquiwrap@0.0.1
```

## List of use cases
1. Get a list of matches with information(picks/bans/meta) for a given tournament
2. Get a list of tournaments (with meta) for a given link from liquipediaa
3. Get current matches (going live) once mode
4. Get current (going live) with the option to use frequently (prevents ban from liquipedia)
## API Reference

#### Get a list of matches for a given tournament

```go
  	cst := liquiwrap.NewGetCurrentStateOfTheTournament("dota2", "BLAST/Slam/3")
	l, err := cst.Action()
	if err != nil {
		fmt.Println(err)
	}
```


#### Get a list of tournaments by liquipedia url

```go
	nt := liquiwrap.NewGetTournamentListByUrl("https://liquipedia.net/dota2/Tier_1_Tournaments")
	
	ts, err := nt.Action()
	if err != nil {
		fmt.Println(err)
	}
```

#### Get current matches (going live) once mode

```go
	gm := liquiwrap.NewGetOngoingMatchesByGame("dota2", 1)
	ts, err := gm.Action()
	if err != nil {
		fmt.Println(err)
	}
```


#### Get current (going live) with the option to use frequently (prevents ban from liquipedia)

```go
	bus := liquiwrap.NewLiquipediaBus(context.Background())
	go func() {
		bus.MustRun()
	}()

	o := liquiwrap.NewGetOngoingMatchesByGameWithBus("dota2", 1, bus)
	newOngoing, err := o.Action()
	if err != nil {
		fmt.Println(err)
	}
```


