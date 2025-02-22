package backend_golang

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func (a *App) StartServer(python string, port int, host string) (string, error) {
	var err error
	if python == "" {
		python, err = GetPython()
	}
	if err != nil {
		return "", err
	}
	return Cmd(python, "./backend-python/main.py", strconv.Itoa(port), host)
}

func (a *App) ConvertModel(python string, modelPath string, strategy string, outPath string) (string, error) {
	var err error
	if python == "" {
		python, err = GetPython()
	}
	if err != nil {
		return "", err
	}
	return Cmd(python, "./backend-python/convert_model.py", "--in", modelPath, "--out", outPath, "--strategy", strategy)
}

func (a *App) DepCheck(python string) error {
	var err error
	if python == "" {
		python, err = GetPython()
	}
	if err != nil {
		return err
	}
	out, err := exec.Command(python, a.exDir+"./backend-python/dep_check.py").CombinedOutput()
	if err != nil {
		return errors.New("DepCheck Error: " + string(out))
	}
	return nil
}

func (a *App) InstallPyDep(python string, cnMirror bool) (string, error) {
	var err error
	if python == "" {
		python, err = GetPython()
		if runtime.GOOS == "windows" {
			python = `"%CD%/` + python + `"`
		}
	}
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		ChangeFileLine("./py310/python310._pth", 3, "Lib\\site-packages")
		installScript := python + " ./backend-python/get-pip.py -i https://pypi.tuna.tsinghua.edu.cn/simple\n" +
			python + " -m pip install torch==1.13.1 torchvision==0.14.1 torchaudio==0.13.1 --index-url https://download.pytorch.org/whl/cu117\n" +
			python + " -m pip install -r ./backend-python/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple\n" +
			"exit"
		if !cnMirror {
			installScript = strings.Replace(installScript, " -i https://pypi.tuna.tsinghua.edu.cn/simple", "", -1)
			installScript = strings.Replace(installScript, "requirements.txt", "requirements_versions.txt", -1)
		}
		err = os.WriteFile("./install-py-dep.bat", []byte(installScript), 0644)
		if err != nil {
			return "", err
		}
		return Cmd("install-py-dep.bat")
	}

	if cnMirror {
		return Cmd(python, "-m", "pip", "install", "-r", "./backend-python/requirements_without_cyac.txt", "-i", "https://pypi.tuna.tsinghua.edu.cn/simple")
	} else {
		return Cmd(python, "-m", "pip", "install", "-r", "./backend-python/requirements_without_cyac.txt")
	}
}
