// datatype.go
package ansible

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"time"
)

// InventoryFile represents the object structure of ansible inventory
/**

 */

var LineBreak = "\n"

type InventoryFile struct {
	sections map[string][]string
	keys     []string
}

func NewInventory() *InventoryFile {

	sections := make(map[string]([]string), 0)

	keys := make([]string, 0)

	return &InventoryFile{
		sections: sections,
		keys:     keys,
	}
}

// Add section to inventory file, it should be on order
func (f *InventoryFile) AddSection(sectionName string, members []string) {

	if _, ok := f.sections[sectionName]; !ok {
		//Add new section keys
		f.keys = append(f.keys, sectionName)
	}

	f.sections[sectionName] = members
}

// SaveTo write the content into file system
func (f *InventoryFile) SaveTo(filename string) error {
	return f.SaveToIndent(filename, "")
}

// Save the inventory file content into the file with given indention
func (f *InventoryFile) SaveToIndent(filename, indent string) error {
	tmpPath := filename + "." + strconv.Itoa(time.Now().Nanosecond()) + ".tmp"
	defer os.Remove(tmpPath)

	fw, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	if _, err = f.WriteToIndent(fw, indent); err != nil {
		fw.Close()
		return err
	}
	fw.Close()

	// Remove old file and rename the new one.
	os.Remove(filename)
	return os.Rename(tmpPath, filename)
}

// WriteToIndent writes content into io.Writer with given indention.
func (f *InventoryFile) WriteToIndent(w io.Writer, indent string) (n int64, err error) {

	buf := bytes.NewBuffer(nil)

	for sName, sValue := range f.sections {

		if _, err := buf.WriteString("[" + sName + "]" + LineBreak); err != nil {
			return 0, err
		}

		for _, value := range sValue {
			if _, err := buf.WriteString(value + LineBreak); err != nil {
				return 0, err
			}
		}

		if _, err := buf.WriteString(LineBreak); err != nil {
			return 0, err
		}

	}

	return buf.WriteTo(w)
}
