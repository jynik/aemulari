package cmdline

import (
	ae "../../aemulari"
)


// Scan arguments for help request. This should be called before Parse() so
// that the help request takes precedence over any invalid flag usages
func HelpRequested(args []string) bool {
	for _, arg := range(args) {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

// Parse command line arguments (argv) based upon list of supported arguments.
// Configures and returns an Architecture and Debugger on success.
// Prints errors to stderr and exits the program on failure.
func Parse(supported SupportedArgs, argv []string) (*ae.Architecture, *ae.Debugger) {
	return nil, nil
}
