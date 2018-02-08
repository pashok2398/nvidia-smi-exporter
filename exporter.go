package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const nvidiaSmiCmd = "nvidia-smi"

func fake_metrics(response http.ResponseWriter) {
	var result []string
	result = append(result, "nvidia_fan_speed{gpu=\"0\", name=\"Tesla K80\"} [Not Supported]")
	result = append(result, "nvidia_temperature_gpu{gpu=\"0\", name=\"Tesla K80\"} 32")
	result = append(result, "nvidia_clocks_gr{gpu=\"0\", name=\"Tesla K80\"} 324")
	result = append(result, "nvidia_clocks_sm{gpu=\"0\", name=\"Tesla K80\"} 324")
	result = append(result, "nvidia_clocks_mem{gpu=\"0\", name=\"Tesla K80\"} 324")
	result = append(result, "nvidia_power_draw{gpu=\"0\", name=\"Tesla K80\"} 25.84")
	result = append(result, "nvidia_utilization_gpu{gpu=\"0\", name=\"Tesla K80\"} 70")
	result = append(result, "nvidia_utilization_memory{gpu=\"0\", name=\"Tesla K80\"} 10")
	result = append(result, "nvidia_memory_total{gpu=\"0\", name=\"Tesla K80\"} 11441")
	result = append(result, "nvidia_memory_free{gpu=\"0\", name=\"Tesla K80\"} 4576")
	result = append(result, "nvidia_memory_used{gpu=\"0\", name=\"Tesla K80\"} 6865")
	
	result = append(result, "nvidia_fan_speed{gpu=\"1\", name=\"Tesla K80\"} [Not Supported]")
	result = append(result, "nvidia_temperature_gpu{gpu=\"1\", name=\"Tesla K80\"} 32")
	result = append(result, "nvidia_clocks_gr{gpu=\"1\", name=\"Tesla K80\"} 324")
	result = append(result, "nvidia_clocks_sm{gpu=\"1\", name=\"Tesla K80\"} 324")
	result = append(result, "nvidia_clocks_mem{gpu=\"1\", name=\"Tesla K80\"} 324")
	result = append(result, "nvidia_power_draw{gpu=\"1\", name=\"Tesla K80\"} 25.84")
	result = append(result, "nvidia_utilization_gpu{gpu=\"1\", name=\"Tesla K80\"} 70")
	result = append(result, "nvidia_utilization_memory{gpu=\"1\", name=\"Tesla K80\"} 10")
	result = append(result, "nvidia_memory_total{gpu=\"1\", name=\"Tesla K80\"} 11441")
	result = append(result, "nvidia_memory_free{gpu=\"1\", name=\"Tesla K80\"} 4576")
	result = append(result, "nvidia_memory_used{gpu=\"1\", name=\"Tesla K80\"} 6865")

	fmt.Fprintf(response, strings.Join(result, "\n"))
}

func metrics(response http.ResponseWriter, request *http.Request) {

	out, err := exec.Command(
		nvidiaSmiCmd,
		"--query-gpu=name,index,fan.speed,temperature.gpu,clocks.gr,clocks.sm,clocks.mem,power.draw,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used",
		"--format=csv,noheader,nounits").Output()

	if err != nil {
		fmt.Printf("%s\n", err)
		fake_metrics(response)
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	metricList := []string{
		"nvidia_fan_speed", "nvidia_temperature_gpu", "nvidia_clocks_gr", "nvidia_clocks_sm", "nvidia_clocks_mem", "nvidia_power_draw",
		"nvidia_utilization_gpu", "nvidia_utilization_memory", "nvidia_memory_total", "nvidia_memory_free", "nvidia_memory_used"}

	result := ""
	for _, row := range records {
		name := fmt.Sprintf("%s", row[0])
		index := fmt.Sprintf("%s", row[1])
		for idx, value := range row[2:] {
			result = fmt.Sprintf(
				"%s%s{gpu=\"%s\", name=\"%s\"} %s\n", result,
				metricList[idx], index, name, value)
		}
	}

	fmt.Fprintf(response, result)
}

func home(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "<html><head><title>Nvidia SMI Exporter</title></head><body><h1>Nvidia SMI Exporter</h1><p><a href=\"/metrics\">Metrics</a></p></body></html>")
}

func main() {
	addr := ":9102"
	if len(os.Args) > 1 {
		addr = ":" + os.Args[1]
	}

	http.HandleFunc("/metrics", metrics)
	http.HandleFunc("/", home)

	log.Printf("Listen port http://localhost%s\n", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}