package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) != 3 {
		panic("expecting a formPath and an outPath")
	}
	pathA := os.Args[1]
	outPath := os.Args[2]
	convertJSONtoAndroidXML(pathA, outPath)
}

func convertJSONtoAndroidXML(formPath, outPath string) error {
	os.MkdirAll(outPath, 0777)

	rawJSON, err := os.ReadFile(formPath)
	if err != nil {
		errors.Wrap(err, "json error")
	}

	formObjects := make([]map[string]string, 0)
	json.Unmarshal(rawJSON, &formObjects)

	xmlStr := `<?xml version="1.0" encoding="utf-8"?>
    <LinearLayout xmlns:android="http://schemas.android.com/apk/res/android"
        xmlns:tools="http://schemas.android.com/tools"
        android:layout_width="match_parent"
        android:layout_height="match_parent"
        android:orientation="vertical"
        android:gravity="center"
        tools:context=".FormsActivity">`

	var stringsXML string
	for _, obj := range formObjects {
		if slices.Index(strings.Split(obj["attributes"], ";"), "hidden") != -1 {
			continue
		}

		var addRequiredStr string
		if slices.Index(strings.Split(obj["attributes"], ";"), "required") != -1 {
			addRequiredStr = `android:tag="required"`
		}
		// add label
		xmlStr += fmt.Sprintf(`<TextView
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="%s"
        android:layout_gravity="center_horizontal"
        android:layout_marginBottom="24dp" />
		`, obj["label"])

		xmlStr += "\n"
		switch obj["fieldtype"] {
		case "int":
			xmlStr += fmt.Sprintf(`<EditText
        android:id="@+id/%s"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Enter a number"
				%s
        android:inputType="number" />`, obj["name"], addRequiredStr)

		case "string":
			xmlStr += fmt.Sprintf(`<EditText
        android:id="@+id/%s"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Enter your text here"
				%s
        android:inputType="text" />`, obj["name"], addRequiredStr)

		case "text":
			xmlStr += fmt.Sprintf(`<EditText
        android:id="@+id/%s"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Enter your text here"
        android:inputType="textMultiLine"
    		android:minLines="3"
				%s
    		android:maxLines="6" />`, obj["name"], addRequiredStr)

		case "email":
			xmlStr += fmt.Sprintf(`<EditText
    android:id="@+id/%s"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    android:hint="Email Address"
    android:inputType="textEmailAddress"
		%s
    android:autofillHints="emailAddress" />`, obj["name"], addRequiredStr)

		case "date":
			xmlStr += fmt.Sprintf(`<EditText
        android:id="@+id/%s"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Select Date"
        android:focusable="false"
        android:clickable="true"
				%s
        android:inputType="none" />`, obj["name"], addRequiredStr)

		case "datetime":
			xmlStr += fmt.Sprintf(`<EditText
        android:id="@+id/%s"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:hint="Select DateTime"
        android:focusable="false"
        android:clickable="true"
				%s
        android:inputType="none" />`, obj["name"], addRequiredStr)

		case "select", "multi_display_select", "single_display_select":
			var spinnerItemsFrag string
			for _, item := range strings.Split(obj["select_options"], "\n") {
				spinnerItemsFrag += fmt.Sprintf("<item>%s</item>", item)
			}

			itemsXML := fmt.Sprintf(`
        <string-array name="%s_items">
				%s
        </string-array>`, obj["name"], spinnerItemsFrag)

			stringsXML += itemsXML + "\n"
			xmlStr += fmt.Sprintf(`<Spinner
        android:id="@+id/%s"
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:entries="@array/%s_items" /> `, obj["name"], obj["name"])

		case "check":
			xmlStr += fmt.Sprintf(`<CheckBox
        android:id="@+id/%s"
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="Check the box" />`, obj["name"])

		}
		xmlStr += "\n"
		xmlStr += "</LinearLayout>"

		// write xml to outPath
		xmlFileName := strings.ReplaceAll(filepath.Base(formPath), ".f8p", ".xml")
		xmlOutPath := filepath.Join(outPath, xmlFileName)
		os.WriteFile(xmlOutPath, []byte(xmlStr), 0777)

		// append to strings xml
		xmlOutPath2 := filepath.Join(outPath, "append_strings.xml")
		os.WriteFile(xmlOutPath2, []byte(stringsXML), 0777)

	}
	return nil
}
