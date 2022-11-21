package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/samwho/streamdeck"
	jxa "github.com/wobsoriano/go-jxa"
)

func main() {
	f, err := ioutil.TempFile("", "com.github.jacobfg.streamdeck-plugin-find-google-meet.log")
	if err != nil {
		log.Fatalf("error creating temp file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	setup(client)

	return client.Run()
}

func setup(client *streamdeck.Client) {

	action := client.Action("com.github.jacobfg.streamdeck-plugin-find-google-meet.action")
	action.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Default().Printf("KeyDown: %+v", event)

		code := `
		(function() {
			var chrome = Application('Google Chrome');
			chrome.activate();
			for (win of chrome.windows()) {
			  var tabIndex =
				win.tabs().findIndex(tab => tab.url().match(/meet.google.com/));
			  if (tabIndex != -1) {
				win.activeTabIndex = (tabIndex + 1);
				win.index = 1;
			  }
			}
		  })();
			`
		v, err := jxa.RunJXA(code)

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Default().Printf("Is dark mode: %s", v)

		return nil
	})
}
