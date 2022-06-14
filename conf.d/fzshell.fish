function _fzshell_install --on-event fzshell_install
    set --query XDG_DATA_HOME || set --local XDG_DATA_HOME ~/.local/share
    set --universal fzshell_data $XDG_DATA_HOME/fzshell

    if test ! -d $fzshell_data
        command mkdir -p $fzshell_data
        echo "Downloading fzshell" 2>/dev/null
        if not command git clone https://github.com/mnowotnik/fzshell $fzshell_data
            echo "fzshell: Can't git clone fzshell project"
            return 1
        end
        cd $fzshell_data
    else
        echo "fzshell already exists. Updating" 2>/dev/null
        cd $fzshell_data
        if not command git pull
            echo "fzshell: Can't git pull fzshell project"
            return 1
        end
    end

    bash scripts/install.sh --no-instructions --install-fish-keys
end

function _fzshell_update --on-event fzshell_update
    cd $fzshell_data
    if not command git pull
        echo "fzshell: Can't git pull fzshell project"
        return 1
    end
    bash scripts/install.sh --no-instructions --install-fish-keys
end

function _fzshell_uninstall --on-event fzshell_uninstall
    command rm -rf $fzshell_data
    set -U fzshell_data
    set --query XDG_CONFIG_HOME || set --local XDG_CONFIG_HOME ~/.config
    if test ! -d $XDG_CONFIG_HOME
        return
    end

    rm -f "$XDG_CONFIG_HOME/fish/conf.d/fzshell_key_bindings.fish"
end
