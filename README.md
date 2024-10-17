# Welcome to BubbleUp

Float your alerts to the top of your TUI like a bubble in a soda. Integrates with BubbleTea applications seamlessly to render your status updates in style.

![Example GIF](./examples/example.gif)

![Keypress Monitoring GIF](./examples/keypresses.gif)

## Getting Started

Run the following to download the module:

```sh
go get go.dalton.dog/bubbleup
```

Then import it into your project with the following:

```go
import (
    "go.dalton.dog/bubbleup"
    tea "github.com/charmbracelet/bubbletea"
)
```

This is assuming you already have BubbleTea installed and available in your project. Go check out their repo [here](https://github.com/charmbracelet/bubbletea) for more information!

From there, it's as simple as creating a new `bubbleup.AlertModel` by calling `bubbleup.NewAlertModel([#], [true|false])`, and embedding the returned model inside of your main model.

The first parameter is the width that you want the alerts to render at.  
*Note: If your message length exceeds the max width, it will wrap*  

The second parameter is whether or not to render the included NerdFont symbols for the included alert types, or just to use normal ASCII character prefixes. Check out Nerd Fonts [here](https://nerdfonts.com).  
*Note: If you override or don't use the included alert types, then this parameter doesn't matter.*

## Integrating Into Your BubbleTea App

### Init

Be sure to return the result of the alert models' `Init()`.  
If you need to also return one or more commands, be sure to use `tea.Batch()` to bundle them together.

### Update

This is where you'll actually spawn the alerts in, which is as easy as calling `NewAlertCmd()` with a key and a message. The formatting and stylings are handled by what's provided in the stored `AlertDefinition` types (more info below).  

Example with the included Info alert type:

Be sure to pass any received messages to the alert model, and appropriately use the return values.  
Reassign your stored alert with the updated alert, and return the given command, either alone or via tea.Batch().  
```go
alertCmd = m.alert.NewAlertCmd(bubbleup.InfoKey, "New info alert.") // Get the command to initiate the desired alert

outAlert, outCmd := m.alert.Update(msg)  // Pass any messages to the alert model, such as alert or tick messages
m.alert = outAlert.(bubbleup.AlertModel) // Reassign the returned alert model to the main model

return m, tea.Batch(alertCmd, outCmd)
```

### View

You want to do all of your normal View stuff to render your output, and THEN pass that into your alert model's `Render()` function. This will overlay the alert onto the provided content. It is recommended to have this be the last thing you do in your `View()` function.

*Note: The AlertModels' View() function is empty and is not intended to be called.*

## Creating Your Own Alert Types

You can create your own alert types by creating an instance of an `AlertDefinition` struct, and passing it into your model's `RegisterNewAlertType()` function. The `AlertDefinition` consists of the following parts:  
- Key: (Required) Unique identifier for your alert type. What is passed into `NewAlertCmd` to get rendering information.
- ForeColor: (Required) A hex color string that you want to use as the foreground color of your alert type.
- Style: (Optional) A `lipgloss.Style` struct that will override the default one, but it's up to you to make sure your override meshes well.
- Prefix: (Optional) The symbol or strings used to prefix your message contents. Can be left empty


### Example

You would declare and register your new alert type like this:

```go
    myCustomAlert := AlertDefinition{
        Key: "CoolAlert",
        ForeColor: "#123456",
        Prefix: ":)"
    }

    m.alertModel.RegisterNewAlertType(myCustomAlert)
```

Note that I did not pass a style in, so it'll use the default.

Then call it later by doing:

```go
outAlertCmd := m.alert.NewAlertCmd("CoolAlert", "My really cool alert message")
```

and handling it all as described in the `Update` section above.
