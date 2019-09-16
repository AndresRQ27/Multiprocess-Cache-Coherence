package main

import (
    "github.com/gotk3/gotk3/gtk"
    "log"
	"math/rand"
	"time"
	"sync"
	"fmt"
	"strconv"
)

//BlocksInMemory - int that has the amount of memory lines in the simulated program
const BlocksInMemory = 16

//BlocksInCache - int that has the amount of cache lines in the simulated program
const BlocksInCache = 8

//Clock - int the sets the clock of the processor in seconds
const Clock = 2 * time.Second

//InstructionCounter - counter of the number of instructions executed
var InstructionCounter = 0

//Mean - bigger = +STR instructions / smaller = +LDR instructions
const Mean = 10
//StdDev - bigger = +STD/LDR instructions / smaller = -STD/LDR instructions
const StdDev = 2
/////Default values generates random distributed numbers between 0 and 20

//Starts the interface
func main() {
	// Initialize GTK without parsing any command line arguments.
    gtk.Init(nil)

	//Constructor for the GUI
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Error:",err)
	}

	//Load the GUI from the Glade file
	err = builder.AddFromFile("../resources/window_main.glade")
	if err != nil {
		log.Fatal("Error:",err)
	}

	//Obtains the object window_main via ID
	obj, err := builder.GetObject("window_main")
	if err != nil {
        log.Fatal("Error:", err)
	}
	
	//Converts the object to a gtk.Window
	//Connects the destroy signal to the MainQuit
	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
        gtk.MainQuit()
	})

    // Recursively show all widgets contained in this window.
	win.ShowAll()

	//Starts the CPU process
	go ComputerStart(builder)

	// Begin executing the GTK main loop.  This blocks until
    // gtk.MainQuit() is run.
    gtk.Main()
}

//ComputerStart - function that instantiate all the necessary components for the processor
func ComputerStart(builder *gtk.Builder) {

	//Shared resources created
	SharedMux := sync.Mutex{}
	SharedMemory := NewMemory()
	GUIChannel := make(chan Message, 24)

	//Max buffer of the channel can hold 3 responses from needed data and 3 broadcast messages
	CPU0PublicChannel := make(chan Message, 6)
	CPU1PublicChannel := make(chan Message, 6)
	CPU2PublicChannel := make(chan Message, 6)
	CPU3PublicChannel := make(chan Message, 6)

	CPU0PrivateChannel := make(chan Message, 6)
	CPU1PrivateChannel := make(chan Message, 6)
	CPU2PrivateChannel := make(chan Message, 6)
	CPU3PrivateChannel := make(chan Message, 6)

	//CPU0 creation
	CPU0ProcessorChannel := make(chan Message)
	CPU0Rand := rand.New(rand.NewSource(0))
	
	CPU0CC := CacheController{
		Name:"CPU0",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU0ProcessorChannel,
		SharedMemory:SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,

		GUIChannel:&GUIChannel,
	}

	CPU0Processor := Processor{
		Name:"CPU0",
		InstructionNumber:0,
		PublicCacheController:&CPU0ProcessorChannel,
		Random:CPU0Rand,
		GUIChannel:&GUIChannel,
	}

	//CPU1 creation
	CPU1ProcessorChannel := make(chan Message)
	CPU1Rand := rand.New(rand.NewSource(1))
	
	CPU1CC := CacheController{
		Name:"CPU1",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU1ProcessorChannel,
		SharedMemory:SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,

		GUIChannel:&GUIChannel,
	}

	CPU1Processor := Processor{
		Name:"CPU1",
		InstructionNumber:0,
		PublicCacheController:&CPU1ProcessorChannel,
		Random:CPU1Rand,
		GUIChannel:&GUIChannel,
	}

	//CPU2 creation
	CPU2ProcessorChannel := make(chan Message)
	CPU2Rand := rand.New(rand.NewSource(2))
	
	CPU2CC := CacheController{
		Name:"CPU2",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU2ProcessorChannel,
		SharedMemory:SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,

		GUIChannel:&GUIChannel,
	}

	CPU2Processor := Processor{
		Name:"CPU2",
		InstructionNumber:0,
		PublicCacheController:&CPU2ProcessorChannel,
		Random:CPU2Rand,
		GUIChannel:&GUIChannel,
	}

	//CPU3 creation
	CPU3ProcessorChannel := make(chan Message)
	CPU3Rand := rand.New(rand.NewSource(3))
	
	CPU3CC := CacheController{
		Name:"CPU3",
		PrivateCache:map[int]*CacheLine{
			0:EmptyCacheLine(),
			1:EmptyCacheLine(),
			2:EmptyCacheLine(),
			3:EmptyCacheLine(),
			4:EmptyCacheLine(),
			5:EmptyCacheLine(),
			6:EmptyCacheLine(),
			7:EmptyCacheLine(),
		},
		PrivateProcessor:&CPU3ProcessorChannel,
		SharedMemory:SharedMemory,
		Mux:&SharedMux,
		PublicChannelCPU0:&CPU0PublicChannel,
		PublicChannelCPU1:&CPU1PublicChannel,
		PublicChannelCPU2:&CPU2PublicChannel,
		PublicChannelCPU3:&CPU3PublicChannel,

		PrivateChannelCPU0:&CPU0PrivateChannel,
		PrivateChannelCPU1:&CPU1PrivateChannel,
		PrivateChannelCPU2:&CPU2PrivateChannel,
		PrivateChannelCPU3:&CPU3PrivateChannel,

		GUIChannel:&GUIChannel,
	}

	CPU3Processor := Processor{
		Name:"CPU3",
		InstructionNumber:0,
		PublicCacheController:&CPU3ProcessorChannel,
		Random:CPU3Rand,
		GUIChannel:&GUIChannel,
	}

	//Initiate all the go routines
	go CPU0CC.Listen()
	go CPU0Processor.ExecuteInstruction()

	go CPU1CC.Listen()
	go CPU1Processor.ExecuteInstruction()

	go CPU2CC.Listen()
	go CPU2Processor.ExecuteInstruction()

	go CPU3CC.Listen()
	go CPU3Processor.ExecuteInstruction()

	go UpdateGUI(builder,&CPU0CC.PrivateCache, &CPU1CC.PrivateCache, 
		&CPU2CC.PrivateCache, &CPU3CC.PrivateCache,&SharedMemory.MemoryMap,&GUIChannel)

	time.Sleep(Clock) //Sleeps for 1 second during the first start

	//Infinite loop that manages the counter
	for {
		InstructionCounter++ //Add 1 to the instruction counter
		time.Sleep(5*Clock) //Sleeps for 3 seconds
	}		

	return
}

//UpdateGUI - function that continously updates the GUI with the current values
func UpdateGUI(builder *gtk.Builder, CPU0, CPU1, CPU2, CPU3 *map[int]*CacheLine, 
	theMemory *map[int]string, GUIChannel *chan Message) {
	guiCounter := 0
	var memoryArray [16]*gtk.Label

	var valueArray0 [8]*gtk.Label
	var valueArray1 [8]*gtk.Label
	var valueArray2 [8]*gtk.Label
	var valueArray3 [8]*gtk.Label

	var tagArray0 [8]*gtk.Label
	var tagArray1 [8]*gtk.Label
	var tagArray2 [8]*gtk.Label
	var tagArray3 [8]*gtk.Label

	var stateArray0 [8]*gtk.Label
	var stateArray1 [8]*gtk.Label
	var stateArray2 [8]*gtk.Label
	var stateArray3 [8]*gtk.Label

	var textLog [4]*gtk.TextBuffer

	//Value for the memoryArray of labels
	obj, err := builder.GetObject("MemValue0")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[0] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue1")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue2")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue3")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue4")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue5")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue6")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue7")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue8")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[8] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue9")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[9] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue10")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[10] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue11")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[11] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue12")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[12] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue13")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[13] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue14")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[14] = obj.(*gtk.Label)
	obj, err = builder.GetObject("MemValue15")
	if err != nil {log.Fatal("Error:", err)}
	memoryArray[15] = obj.(*gtk.Label)

	//Values for the valueArray0 of labels
	obj, err = builder.GetObject("CacheValue01")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue02")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue03")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue04")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue05")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue06")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue07")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue00")
	if err != nil {log.Fatal("Error:", err)}
	valueArray0[0] = obj.(*gtk.Label)

	//Values for the valueArray1 of labels
	obj, err = builder.GetObject("CacheValue11")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue12")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue13")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue14")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue15")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue16")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue17")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue10")
	if err != nil {log.Fatal("Error:", err)}
	valueArray1[0] = obj.(*gtk.Label)

	//Values for the valueArray2 of labels
	obj, err = builder.GetObject("CacheValue21")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue22")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue23")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue24")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue25")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue26")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue27")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue20")
	if err != nil {log.Fatal("Error:", err)}
	valueArray2[0] = obj.(*gtk.Label)

	//Values for the valueArray3 of labels
	obj, err = builder.GetObject("CacheValue31")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue32")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue33")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue34")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue35")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue36")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue37")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("CacheValue30")
	if err != nil {log.Fatal("Error:", err)}
	valueArray3[0] = obj.(*gtk.Label)

	//Values of the textLogs of textbuffer
	obj, err = builder.GetObject("CPU0Log")
	if err != nil {log.Fatal("Error:", err)}
	textLog[0] = obj.(*gtk.TextBuffer)
	obj, err = builder.GetObject("CPU1Log")
	if err != nil {log.Fatal("Error:", err)}
	textLog[1] = obj.(*gtk.TextBuffer)
	obj, err = builder.GetObject("CPU2Log")
	if err != nil {log.Fatal("Error:", err)}
	textLog[2] = obj.(*gtk.TextBuffer)
	obj, err = builder.GetObject("CPU3Log")
	if err != nil {log.Fatal("Error:", err)}
	textLog[3] = obj.(*gtk.TextBuffer)




	//Values for the tagArray0 of labels
	obj, err = builder.GetObject("Tag01")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag02")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag03")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag04")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag05")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag06")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag07")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag00")
	if err != nil {log.Fatal("Error:", err)}
	tagArray0[0] = obj.(*gtk.Label)

	//Values for the valueArray1 of labels
	obj, err = builder.GetObject("Tag11")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag12")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag13")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag14")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag15")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag16")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag17")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag10")
	if err != nil {log.Fatal("Error:", err)}
	tagArray1[0] = obj.(*gtk.Label)

	//tags for the tagArray2 of labels
	obj, err = builder.GetObject("Tag21")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag22")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag23")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag24")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag25")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag26")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag27")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag20")
	if err != nil {log.Fatal("Error:", err)}
	tagArray2[0] = obj.(*gtk.Label)

	//tags for the tagArray3 of labels
	obj, err = builder.GetObject("Tag31")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag32")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag33")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag34")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag35")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag36")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag37")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("Tag30")
	if err != nil {log.Fatal("Error:", err)}
	tagArray3[0] = obj.(*gtk.Label)




	//Values for the tagArray0 of labels
	obj, err = builder.GetObject("State01")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State02")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State03")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State04")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State05")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State06")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State07")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State00")
	if err != nil {log.Fatal("Error:", err)}
	stateArray0[0] = obj.(*gtk.Label)

	//Values for the valueArray1 of labels
	obj, err = builder.GetObject("State11")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State12")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State13")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State14")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State15")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State16")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State17")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State10")
	if err != nil {log.Fatal("Error:", err)}
	stateArray1[0] = obj.(*gtk.Label)

	//states for the stateArray2 of labels
	obj, err = builder.GetObject("State21")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State22")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State23")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State24")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State25")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State26")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State27")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State20")
	if err != nil {log.Fatal("Error:", err)}
	stateArray2[0] = obj.(*gtk.Label)

	//states for the stateArray3 of labels
	obj, err = builder.GetObject("State31")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[1] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State32")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[2] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State33")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[3] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State34")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[4] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State35")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[5] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State36")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[6] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State37")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[7] = obj.(*gtk.Label)
	obj, err = builder.GetObject("State30")
	if err != nil {log.Fatal("Error:", err)}
	stateArray3[0] = obj.(*gtk.Label)

	for {
		if (InstructionCounter - guiCounter) > 0 {
			for i := range memoryArray {
				memoryArray[i].SetLabel((*theMemory)[i])
			}
			for i := range valueArray0 {
				valueArray0[i].SetLabel((*CPU0)[i].Data)
				valueArray1[i].SetLabel((*CPU1)[i].Data)
				valueArray2[i].SetLabel((*CPU2)[i].Data)
				valueArray3[i].SetLabel((*CPU3)[i].Data)

				tagArray0[i].SetLabel(strconv.Itoa((*CPU0)[i].tag))
				tagArray1[i].SetLabel(strconv.Itoa((*CPU1)[i].tag))
				tagArray2[i].SetLabel(strconv.Itoa((*CPU2)[i].tag))
				tagArray3[i].SetLabel(strconv.Itoa((*CPU3)[i].tag))

				stateArray0[i].SetLabel((*CPU0)[i].State)
				stateArray1[i].SetLabel((*CPU1)[i].State)
				stateArray2[i].SetLabel((*CPU2)[i].State)
				stateArray3[i].SetLabel((*CPU3)[i].State)
			}
			for len(*GUIChannel) > 0 {
				switch msg := <-*GUIChannel; msg.CPU {
				case "CPU0":
					end := textLog[0].GetEndIter()
					textLog[0].Insert(end,msg.Value)
				case "CPU1":
					end := textLog[1].GetEndIter()
					textLog[1].Insert(end,msg.Value)
				case "CPU2":
					end := textLog[2].GetEndIter()
					textLog[2].Insert(end,msg.Value)
				case "CPU3":
					end := textLog[3].GetEndIter()
					textLog[3].Insert(end,msg.Value)
				}
			}
			fmt.Println("Clock",guiCounter)
			guiCounter = InstructionCounter
		}
	}
}