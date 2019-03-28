package main

import (
	"fmt"
	"math/rand"
	"time"
)

// init
// 1. [fq->bam]
// 2. [bam->mergedBam]
// 3. [mergedBam->realign]
// 4. [realignBam->BQSR]
// 5. [BQSR->final.bam]
// 6 [bam->qc],[bam->CNV]
// 5 [BQSR->Vcf]
// 6 [vcf->final.vcf]
// 7 [final.vcf->anno]

type link struct {
	From []chan string
	To   []chan string
}

func main() {
	sampleList := []string{"sample1", "sample2"}

	var start = make(chan string)
	var startSample []chan string
	for range sampleList {
		startSample = append(startSample, make(chan string))
	}
	startSampleLink := link{
		[]chan string{start},
		startSample,
	}
	runTask(fmt.Sprintf("%s", "start samples"), startSampleLink)

	var allBam []chan string
	var allQc []chan string
	var allAnno []chan string
	var excelFrom []chan string

	for i, sample := range sampleList {
		lanes := []string{"lane1", "lane2"}

		var bwaTo []chan string
		for range lanes {
			bwaTo = append(bwaTo, make(chan string))
		}
		bwaLink := link{
			[]chan string{startSample[i]},
			bwaTo,
		}

		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "start bwa"), bwaLink)
		var mergeFrom []chan string
		for j, fq := range lanes {
			to := make(chan string)
			bwaLaneLink := link{
				[]chan string{bwaTo[j]},
				[]chan string{to},
			}
			message := fmt.Sprintf("%d-%d:%s-%s\t%s", i, j, sample, fq, "fq->bam")
			runTask(message, bwaLaneLink)
			mergeFrom = append(mergeFrom, bwaLaneLink.To...)
		}

		to := make(chan string)
		mergeLink := link{
			mergeFrom,
			[]chan string{to},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "mergeBam"), mergeLink)

		to = make(chan string)
		realignLink := link{
			mergeLink.To,
			[]chan string{to},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "realign"), realignLink)

		bqsrLink := link{
			realignLink.To,
			[]chan string{make(chan string), make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "BQSR"), bqsrLink)

		finalBamLink := link{
			[]chan string{bqsrLink.To[0]},
			[]chan string{make(chan string), make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "printReads"), finalBamLink)

		bamQcLink := link{
			[]chan string{finalBamLink.To[0]},
			[]chan string{make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "bamQC"), bamQcLink)
		allQc = append(allQc, bamQcLink.To[0])
		allBam = append(allBam, finalBamLink.To[1])

		vcfLink := link{
			[]chan string{bqsrLink.To[1]},
			[]chan string{make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "call VCF"), vcfLink)

		finalVcfLink := link{
			[]chan string{vcfLink.To[0]},
			[]chan string{make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "final VCF"), finalVcfLink)

		annoLink := link{
			[]chan string{finalVcfLink.To[0]},
			[]chan string{make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "anno VCF"), annoLink)
		allAnno = append(allAnno, annoLink.To[0])
	}

	CNVLink := link{
		allBam,
		[]chan string{make(chan string)},
	}
	runTask(fmt.Sprintf("%d:%s\t%s", "batch", "all", "CNV"), CNVLink)

	excelFrom = append(excelFrom, allAnno...)
	excelFrom = append(excelFrom, CNVLink.To[0])
	excelFrom = append(excelFrom, allQc...)

	var allExcelFrom []chan string
	for range sampleList {
		allExcelFrom = append(allExcelFrom, make(chan string))
	}
	excelLink := link{
		excelFrom,
		allExcelFrom,
	}
	runTask(fmt.Sprintf("%s:%s\t%s", "each", "samples", "Excel"), excelLink)

	var allExcelDone []chan string
	for i, sample := range sampleList {
		toExcelLink := link{
			[]chan string{excelLink.To[i]},
			[]chan string{make(chan string)},
		}
		runTask(fmt.Sprintf("%d:%s\t%s", i, sample, "Excel"), toExcelLink)
		allExcelDone = append(allExcelDone, toExcelLink.To[0])
	}

	doneLink := link{
		allExcelDone,
		[]chan string{make(chan string)},
	}
	runTask(fmt.Sprintf("%d:%s\t%s", "batch", "all", "Done"), doneLink)

	start <- fmt.Sprintf("%s", "init")

	fmt.Println(<-doneLink.To[0])
}

func runTask(message string, Link link) {
	from := Link.From
	to := Link.To
	go func() {
		t := rand.Intn(1e1)
		for _, c := range from {
			fmt.Printf("%s done,start %s,took %ds\n%+v\n", <-c, message, t, Link)
		}
		time.Sleep(time.Duration(t) * time.Second)
		for i := range to {
			to[i] <- fmt.Sprintf("%s\t%s:%ds\ttask%d", time.Now().String(), message, t, i)
		}
	}()
	return
}
