package dependencies

import "fmt"

type Dependencies struct {
	dependent map[string][]string
}

var dependencies Dependencies

func NextBatch(deps *Dependencies) (string, error) {
	for key, vals := range deps.dependent {
		if len(vals) == 0 {
			return key, nil
		}
	}
	return "", fmt.Errorf("No dependencies or circular dependencies detected")
}

//func get(path string) (error, deps *Dependencies) {
//  depTree := make(map[string][]string)
//  err := filepath.Walk(path, func(taskPath string, info *os.FileInfo, err error) error {

//    task := filepath.Base(taskPath)

//    // Only run for each of "path/stage" subdirectories
//    if !info.IsDir() || filepath.Base(path) == task {
//      return nil
//    }

//    file, err := os.Open(filepath.Join(taskPath, ".dependencies"))
//    if err != nil {
//      if os.IsNotExist(err) {
//        return nil
//      } else {
//        return err
//      }
//    }

//    scanner := bufio.NewScanner(file)
//    for scanner.Scan() {
//      if scanner.Text != "" {
//        deps[task] = append(deps[task], scanner.Text...)
//      }
//    }

//  })(path)
//  if err != nil {
//    return err
//  }

//  deps := Dependencies{}

//  return
//}
