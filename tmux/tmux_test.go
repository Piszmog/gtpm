package tmux_test

import (
	"os"
	"testing"

	"github.com/Piszmog/gtpm/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPlugins(t *testing.T) {
	tests := []struct {
		name            string
		rawFile         string
		expectedPlugins []string
		err             error
	}{
		{
			name:            "Empty File",
			rawFile:         ``,
			expectedPlugins: nil,
			err:             nil,
		},
		{
			name:            "No Plugins",
			rawFile:         `set -g default-terminal "xterm-256color"`,
			expectedPlugins: nil,
			err:             nil,
		},
		{
			name: "One Plugin",
			rawFile: `
set -g default-terminal "xterm-256color"
set -g @plugin 'tmux-plugins/tmux-sensible'
`,
			expectedPlugins: []string{"tmux-plugins/tmux-sensible"},
			err:             nil,
		},
		{
			name: "One Plugin and another commented",
			rawFile: `
set -g default-terminal "xterm-256color"
set -g @plugin 'tmux-plugins/tmux-sensible'
#set -g @plugin 'odedlaz/tmux-onedark-theme'
`,
			expectedPlugins: []string{"tmux-plugins/tmux-sensible"},
			err:             nil,
		},
		{
			name: "Multiple Plugins",
			rawFile: `
set -g default-terminal "xterm-256color"
set -g @plugin 'tmux-plugins/tmux-sensible'

set -g @plugin 'odedlaz/tmux-onedark-theme'
`,
			expectedPlugins: []string{
				"tmux-plugins/tmux-sensible",
				"odedlaz/tmux-onedark-theme",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.CreateTemp("", "config-*.conf")
			require.NoError(t, err)
			defer os.Remove(f.Name())

			f.Write([]byte(test.rawFile))

			err = f.Close()
			require.NoError(t, err)

			plugins, err := tmux.GetPlugins(f.Name())
			if test.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				assert.Nil(t, plugins)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedPlugins, plugins)
			}
		})
	}
}
