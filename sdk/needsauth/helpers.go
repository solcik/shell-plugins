package needsauth

import (
	"github.com/1Password/shell-plugins/sdk"
)

// IfAll returns a NeedsAuthentication that opts in to the authentication requirement only if
// all the specified rules opt in to the authentication requirement.
func IfAll(rules ...sdk.NeedsAuthentication) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, rule := range rules {
			if !rule(in) {
				return false
			}
		}
		return true
	}
}

// IfAny returns a NeedsAuthentication rule that only opts in to the authentication requirement
// if at least one specified rule opts in to the authentication requirement.
func IfAny(rules ...sdk.NeedsAuthentication) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, rule := range rules {
			if rule(in) {
				return true
			}
		}
		return false
	}
}

// ForCommand returns a NeedsAuthentication rule to require authentication for
// certain (sub)command, e.g. ["account"] or ["account", "list"].
func ForCommand(command ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		if len(command) > len(in.CommandArgs) {
			return false
		}

		for i := range command {
			if command[i] != in.CommandArgs[i] {
				return false
			}
			if i == len(command)-1 {
				return true
			}
		}

		return false
	}
}

// Always returns a NeedsAuthentication rule to always require authentication.
func Always() sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		return true
	}
}

// argsBeforeSeparator returns args up to (but not including) the first "--" separator.
func argsBeforeSeparator(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[:i]
		}
	}
	return args
}

// NotForExactArgs returns a NeedsAuthentication rule to opt out of authentication when
// the command-line args are an exact match with the passed in args.
func NotForExactArgs(argsToSkip ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		args := argsBeforeSeparator(in.CommandArgs)
		if len(args) != len(argsToSkip) {
			return true
		}

		for i, commandArg := range args {
			if commandArg != argsToSkip[i] {
				return true
			}
		}

		return false
	}
}

// NotWhenContainsArgs returns a NeedsAuthentication rule to not require authentication when
// the exact sequence of argsToSkip is present somewhere in the command-line args.
func NotWhenContainsArgs(argsSequence ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		args := argsBeforeSeparator(in.CommandArgs)
		if len(argsSequence) == 0 {
			return true
		}

		if len(argsSequence) > len(args) {
			return true
		}

		for i := range args {
			if i+len(argsSequence) > len(args) {
				return true
			}

			matches := true
			for i, argsToCompare := range args[i : i+len(argsSequence)] {
				if argsToCompare != argsSequence[i] {
					matches = false
				}
			}

			// If the argsToSkip are found in the command-line args, return that the command
			// does not not require authentication
			if matches {
				return false
			}
		}
		return true
	}
}

func NotForHelp() sdk.NeedsAuthentication {
	return IfAll(
		NotWhenContainsArgs("-h"),
		NotWhenContainsArgs("--help"),
		NotWhenContainsArgs("-help"),
		NotWhenContainsArgs("help"),
	)
}

func NotForVersion() sdk.NeedsAuthentication {
	return IfAll(
		NotForExactArgs("-v"),
		NotForExactArgs("--version"),
		NotForExactArgs("-version"),
		NotForExactArgs("version"),
		NotForExactArgs("-V"),
	)
}

func NotWithoutArgs() sdk.NeedsAuthentication {
	return NotForExactArgs()
}

func NotForHelpOrVersion() sdk.NeedsAuthentication {
	return IfAll(NotForHelp(), NotForVersion())
}
