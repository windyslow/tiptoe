package protocol

import (
  "fmt"
  "time"
  "bufio"
  "strconv"
  "os"
  "testing"
  "sync"
)





import (
	"github.com/henrycg/simplepir/pir"
	"github.com/henrycg/simplepir/matrix"
)



import (
  "github.com/fatih/color"
)



func TestCoooooF(t *testing.T) {
	// ip config 
	
	coordinatorAddr := "219.245.186.51:1237"


	//tiptoe database
	color.Yellow("Setting up client...")

	c := NewClient(true /* use coordinator */)
	//c.printStep("Getting metadata")
	//hint := c.getHint(true /* keep conn */, coordinatorAddr)
	//c.Setup(hint)
	//logHintSize(hint)
	//var wg sync.WaitGroup
	
	
	var clients [100]*Client
	var emb_querys	[100]*pir.Query[matrix.Elem64]
	for i:=0;i<len(clients);i++{
		clients[i] = NewClient(true /* use coordinator */)
	}
	for i:=0;i<len(clients);i++ {
		//wg.Add(1)
		emb_querys[i] = clients[i].cycle(coordinatorAddr, true /* verbose */, true /* keep conn */)
	}
	startTime := time.Now()
	var wg sync.WaitGroup	
	for i:=0;i<len(clients);i++{
		wg.Add(1)
		go clients[i].runRound001(coordinatorAddr,emb_querys[i],&wg)
	}
	//wg.Wait()
	duration := time.Since(startTime)
	fmt.Printf("程序运行时间: %dus\n", int(duration.Microseconds()))
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}

func (c *Client) cycle(coordinatorAddr string, verbose, keepConn bool)*pir.Query[matrix.Elem64]{
	//defer wg.Done()
	hint := c.getHint(true /* keep conn */, coordinatorAddr)
	c.printStep("Getting metadata")
	c.Setup(hint)
	logHintSize(hint)
	c.printStep("Running client preprocessing")
	perf := c.preprocessRound(coordinatorAddr, true /* verbose */, true /* keep conn */)
	fmt.Printf("\n\n")
	//c.runRound0(cols,cases,facts,perf, in, out, text, coordinatorAddr, true /* verbose */, true /* keep conn */)
	return c.runRound000(perf, coordinatorAddr, true /* verbose */, true /* keep conn */)
}

func (c *Client) runRound000(p Perf, 
	coordinatorAddr string, verbose, keepConn bool)*pir.Query[matrix.Elem64]{
	//y := color.New(color.FgYellow, color.Bold)

	// Build embeddings query
	start := time.Now()
	if verbose {
	c.printStep("Generating embedding of the query")
	}

	var query struct {
	Cluster_index uint64
	Emb           []int8
	}
	query.Emb = readintegers0("/home/yyh/yyh/similar/tiptoe/search/protocol/data.txt")
	query.Cluster_index = 0

	embQuery := c.QueryEmbeddings(query.Emb, query.Cluster_index)
	print("\ncluster_index:",query.Cluster_index,"\n")
	p.clientSetup = time.Since(start).Seconds()
	// Send embeddings query to server
	if verbose {
	c.printStep("Sending query to server")
	}
	return embQuery
}

func (c *Client) runRound001(coordinatorAddr string, embQuery *pir.Query[matrix.Elem64], wg *sync.WaitGroup)*pir.Answer[matrix.Elem64]{
	defer wg.Done()
	embAns := c.getEmbeddingsAnswer(embQuery, true /* keep conn */, coordinatorAddr)
	return embAns
}


func (c *Client) runRound00(p Perf, 
	coordinatorAddr string, verbose, keepConn bool){
	/*
	//y := color.New(color.FgYellow, color.Bold)

	// Build embeddings query
	start := time.Now()
	if verbose {
	c.printStep("Generating embedding of the query")
	}

	var query struct {
	Cluster_index uint64
	Emb           []int8
	}
	query.Emb = readintegers0("/home/nsklab/yyh/similar/tiptoe/search/mydata/data/image_bytes/test_embedding.txt")
	query.Cluster_index = 0

	embQuery := c.QueryEmbeddings(query.Emb, query.Cluster_index)
	print("\ncluster_index:",query.Cluster_index,"\n")
	p.clientSetup = time.Since(start).Seconds()
	// Send embeddings query to server
	if verbose {
	c.printStep("Sending query to server")
	}
	networkingStart := time.Now()

*/
	//embAns := c.getEmbeddingsAnswer(embQuery, true /* keep conn */, coordinatorAddr)

/*
	p.t1, p.up1, p.down1 = logStats(c.params.NumDocs, networkingStart, embQuery, embAns)

	// Recover document and URL chunk to query for
	c.printStep("Decrypting server answer")
	embDec := c.ReconstructEmbeddingsWithinCluster(embAns, query.Cluster_index)
	scores := embeddings.SmoothResults(embDec, c.embInfo.P())

	docIndex := uint64(maxIndex(scores))
	// build query
	var client_state []mypir.State
	var query0 mypir.MsgSlice
	cs, q := client_pir.Query(uint64(docIndex), offline_data.Shared_state, offline_data.P, offline_data.Info)
	client_state = append(client_state, cs)
	query0.Data = append(query0.Data, q)
	conn, err := net.Dial("tcp", ip_addr)
	if err != nil {
	fmt.Println("Error connecting:", err.Error())
	return
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(query0)
	if err != nil {
	fmt.Println("Error encoding query:", err.Error())
	return
	}

	//log out
	log.Printf("[%v][%v][1. Send built query]\t Elapsed:%v \tSize:%vKB", program, conn.LocalAddr(),
	mypir.PrintTime(start), float64(query0.Size()*uint64(offline_data.P.Logq)/(8.0*1024.0)))
	//receive answer
	var answer mypir.Msg
	decoder = json.NewDecoder(conn)
	err = decoder.Decode(&answer)
	if err != nil {
	fmt.Println("Error decoding answer:", err.Error())
	return
	}
	//log out
	log.Printf("[%v][%v][2. Receive answer]\t Elapsed:%v \tSize:%vKB", program, conn.LocalAddr(),
	mypir.PrintTime(start), float64(answer.Size()*uint64(offline_data.P.Logq)/(8.0*1024.0)))

	//resconstruction
	val := client_pir.StrRecover(uint64(docIndex), uint64(0), offline_data.Offline_download,
	query0.Data[0], answer, offline_data.Shared_state,
	client_state[0], offline_data.P, offline_data.Info)
	val = val
	//print("\ndocindex:",docIndex,"\n")
	//docIndex = uint64(cols[query.Cluster_index][docIndex]%55192)
	//fact := facts[docIndex]
	//docIndex = uint64(cases[docIndex])
	//if verbose {
	//y.Printf("\tDoc %d has the largest inner product with our query\n",
	//docIndex)
	//print(fact,"\n")
	//}
	*/

}

func readintegers0(filename string) []int8{
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close() 
  
	var numbers []int8 
  
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // 逐行读取
		line := scanner.Text() // 获取行文本
		number, err := strconv.Atoi(line) // 将行文本转换为整数
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			continue // 遇到错误时跳过这一行
		}
		numbers = append(numbers, int8(number)) // 将整数添加到切片中
	}
  
	// 检查是否有扫描时发生的错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
  
	return numbers // 打印所有读取的整数
}
