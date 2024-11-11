package main

import (
	"fmt"
	"plusz-backend/scrapper"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	scheduleRevision := scrapper.Scrap(
		"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/zajecia_xml.php?kierunek=IiE&rok=3z",
	)

	fmt.Println(scheduleRevision.Date)
	for _, class := range scheduleRevision.Classes {
		fmt.Println(class)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
