package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	constants "github.com/doezaza12/dynamic-templates/constant"
	"github.com/doezaza12/dynamic-templates/util"
)

func GitClone(url string, noCache bool) error {
	if url == "" {
		return errors.New("url params must not empty")
	}

	if !util.IsRemoteTemplate(url) {
		return errors.New("url is not git's scheme")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	templateDir := filepath.Join(homeDir, constants.APP_STAGING_DIR, constants.DEFAULT_TEMPLATE_DIR)
	if err := os.MkdirAll(templateDir, constants.DEFAULT_FILE_PERMISSION); err != nil {
		if os.IsExist(err) {
			fmt.Printf("%v already existed", templateDir)
		} else {
			return err
		}
	}

	var branchFlag string
	if util.HasRevision(url) {
		baseUrl, revision, found := strings.Cut(url, constants.REVISION_PATTERN)
		if found {
			url = baseUrl
			branchFlag = fmt.Sprintf("-b %v", revision)
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

	cmd := exec.Command(constants.DEFAULT_SHELL, "-c", fmt.Sprintf("git clone --depth=1 %v %v && rm -rf %v/.git", branchFlag, url,
		filepath.Join(templateDir, repoName)))

	cmd.Dir = templateDir

	result, err := cmd.Output()
	if err != nil {
		return err
	}

	log.Println(result)

	return nil
}
