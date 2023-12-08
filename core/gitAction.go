package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GitClone(url string, noCache bool) error {
	if url == "" {
		return errors.New("url params must not empty")
	}

	if !strings.Contains(url, "git@") && !strings.Contains(url, "https") {
		return errors.New("url is not git's scheme")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	templateDir := filepath.Join(homeDir, APP_STAGING_DIR, DEFAULT_TEMPLATE_DIR)

	if err := os.MkdirAll(templateDir, DEFAULT_FILE_PERMISSION); err != nil {
		if os.IsExist(err) {
			fmt.Printf("%v already existed", templateDir)
		} else {
			return err
		}
	}

	repoName := strings.ReplaceAll(filepath.Base(url), ".git", "")
	templateAbsPath := filepath.Join(templateDir, repoName)

	// check if remote template is in local
	fileStat, err := os.Stat(templateAbsPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// if use no cache, it will delete existing template in local if it found one
	if !os.IsNotExist(err) && fileStat.IsDir() && fileStat.Name() == repoName && noCache {
		err := os.RemoveAll(templateAbsPath)
		if err != nil {
			return err
		}
	}

	//  if use cache and template is exist it will bypass git clone
	if !os.IsNotExist(err) && !noCache {
		return nil
	}

	cmd := exec.Command(DEFAULT_SHELL, "-c", fmt.Sprintf("git clone --depth=1 %v && rm -rf %v/.git", url,
		filepath.Join(templateDir, repoName)))

	cmd.Dir = templateDir

	result, err := cmd.Output()
	if err != nil {
		return err
	}

	log.Println(result)

	return nil
}
