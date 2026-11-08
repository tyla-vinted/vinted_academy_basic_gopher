# Go Arena

A shared 2D room. Everyone in the workshop walks around as their own ASCII-art character, talks in speech bubbles, emotes, and high-fives. You write the code that makes your avatar work. At the end, everyone joins one lobby.

## Before the workshop

1. Install Go 1.23 or newer: <https://go.dev/dl/>.
2. Clone this repository.
3. Run the check:

   ```sh
   make doctor
   ```

   Expected result: `All good. Run: make run`.

4. Run the arena in solo mode:

   ```sh
   make run
   ```

   Expected result: a window opens with a grey `[?]`. That is you. You cannot walk yet: that is your first task.

On Linux, Ebitengine needs some system packages. See <https://ebitengine.org/en/documents/install.html>.

On Windows without `make`, use the commands from the [Makefile](Makefile) directly, for example `go run ./cmd`.

## How the workshop works

All your code is in the [`avatar/`](avatar) folder. The rest of the repository (rendering, network) is ready, and you do not have to read it.

You start with step 1 only. Each step is a file with `// TODO`s and a test file. For each step:

1. Write the code in the step file. Run `make run` to see the result in the arena.
2. Run `make check`. Expected result: `Step N is done`. If a test fails, the output tells you what to fix.
3. Run `make next`. Expected result: the next step's files appear in `avatar/`.

`make next` runs the tests first, and it unlocks the next step only when they pass.

Do not edit `avatar/locked.go`. It holds the steps that you have not unlocked yet, and `make next` rewrites it.

| Step | File | Topic | What you see in the arena |
|---|---|---|---|
| 1 | `step1_basics.go` | Functions, multiple return values | Your name and your ASCII art |
| 2 | `step2_movement.go` | `switch` | You walk (WASD or arrow keys) |
| 3 | `step3_looks.go` | `for` loops, `switch` | Your art in your colour, facing the way you walk |
| 4 | `step4_commands.go` | `if`, strings | Chat commands like `/wave` |
| 5 | `step5_structs.go` | Structs and methods | Name cards; high-fives |
| 6 | `step6_pointers.go` | Pointers vs values | `/status` shows under your name; emote cooldown |
| 7 | `step7_generics.go` | Generics | Chat history; the "Nearby" list |
| 8 | `step8_errors.go` | Errors: `error` values, `fmt.Errorf`, `errors.Is`/`errors.As` | Clear error messages in the game and in the terminal |
| 9 | `step9_testing.go` | Testing | Correct greetings on name cards |

## Commands

| Command | What it does |
|---|---|
| `make run` | Opens the arena in solo mode: only you, no server. Chat, emotes and status work. |
| `make check` | Tests the steps that you have, and tells you what to do next. |
| `make next` | If your steps pass, unlocks the next step. |
| `make test` | Runs the tests, with less text than `make check`. |
| `make online SERVER=<url>` | Joins the workshop lobby without the tests. Use it to try high-fives with your neighbours. |
| `make join SERVER=<url>` | When all 9 steps pass, joins the workshop lobby for the finale. |

## Controls

| Key | Action |
|---|---|
| WASD or arrow keys | Move |
| Enter | Chat. Press Enter again to send, Esc to cancel. |
| `/` | Type a command |
| H | High-five the nearest person |
| 1–5 | Emotes: wave, dance, confetti, heart, think |

Chat commands: `/wave`, `/dance`, `/confetti`, `/heart`, `/think`, `/status <text>`, `/status` (clears it), `/highfive`.

## Your character

Your name and your art come from `Me()` in `avatar/step1_basics.go`. The art is at most 5 rows of 10 characters, printable ASCII only, facing right. Write each row in backticks, so that `\` needs no escape. Your colour is the `Color` constant in step 3.

Everyone in the lobby sees your character as **your** code draws it: `NormalizeArt` cleans it, `ColorFor` colours it, and `Mirror` turns it when you walk left.
