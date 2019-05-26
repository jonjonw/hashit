package processor

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	mmapgo "github.com/edsrzf/mmap-go"
	"github.com/minio/blake2b-simd"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
)

func fileProcessorWorker(input chan string, output chan Result) {
	for res := range input {
		if Debug {
			printDebug(fmt.Sprintf("processing %s", res))
		}

		// Open the file and determine if we should read it from disk or memory map
		file, err := os.OpenFile(res, os.O_RDONLY, 0644)

		if err != nil {
			printError(fmt.Sprintf("Unable to process file %s with error %s", res, err.Error()))
			continue
		}

		fi, err := file.Stat()

		if err != nil {
			printError(fmt.Sprintf("Unable to get file info for file %s with error %s", res, err.Error()))
			continue
		}

		fsize := fi.Size()
		_ = file.Close()

		if fsize > StreamSize {
			// If Windows always ignore memory maps and stream the file off disk
			if runtime.GOOS == "windows" || NoMmap == true {

				if Debug {
					printDebug(fmt.Sprintf("%s bytes=%d using scanner", res, fsize))
				}

				// TODO should return a struct with the values we have
				processScanner(res)
			} else {
				if Debug {
					printDebug(fmt.Sprintf("%s bytes=%d using memory map", res, fsize))
				}

				fileStartTime := makeTimestampMilli()
				r, err := processMemoryMap(res)
				if Trace {
					printTrace(fmt.Sprintf("milliseconds processMemoryMap: %s: %d", res, makeTimestampMilli()-fileStartTime))
				}

				if err == nil {
					r.File = res
					r.Bytes = fsize
					output <- r
				}
			}
		} else {
			if Debug {
				printDebug(fmt.Sprintf("%s bytes=%d using read file", res, fsize))
			}

			fileStartTime := makeTimestampNano()
			r, err := processReadFile(res)
			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processReadFile: %s: %d", res, makeTimestampNano()-fileStartTime))
			}

			if err == nil {
				r.File = res
				r.Bytes = fsize
				output <- r
			}
		}
	}
	close(output)
}

// TODO compare this to memory maps
// Random tests indicate that mmap is faster when not in power save mode
func processScanner(filename string) {
	//file, err := os.Open(filename)
	//if err != nil {
	//	printError(fmt.Sprintf("opening file %s: %s", filename, err.Error()))
	//	return
	//}
	//defer file.Close()
	//
	//// Create channels for each hash
	//md5_d := md5.New()
	//sha1_d := sha1.New()
	//sha256_d := sha256.New()
	//sha512_d := sha512.New()
	//blake2b_256_d := blake2b.New256()
	//blake2b_512_d := blake2b.New512()
	//
	//md5c := make(chan []byte, 10)
	//sha1c := make(chan []byte, 10)
	//sha256c := make(chan []byte, 10)
	//sha512c := make(chan []byte, 10)
	//blake2b_256_c := make(chan []byte, 10)
	//blake2b_512_c := make(chan []byte, 10)
	//
	//var wg sync.WaitGroup
	//
	//if hasHash(HashNames.MD5) {
	//	wg.Add(1)
	//	go func() {
	//		for b := range md5c {
	//			md5_d.Write(b)
	//		}
	//		wg.Done()
	//	}()
	//}
	//
	//if hasHash(HashNames.SHA1) {
	//	wg.Add(1)
	//	go func() {
	//		for b := range sha1c {
	//			sha1_d.Write(b)
	//		}
	//		wg.Done()
	//	}()
	//}
	//
	//if hasHash(s_sha256) {
	//	wg.Add(1)
	//	go func() {
	//		for b := range sha256c {
	//			sha256_d.Write(b)
	//		}
	//		wg.Done()
	//	}()
	//}
	//
	//if hasHash(s_sha512) {
	//	wg.Add(1)
	//	go func() {
	//		for b := range sha512c {
	//			sha512_d.Write(b)
	//		}
	//		wg.Done()
	//	}()
	//}
	//
	//if hasHash(s_blake2b256) {
	//	wg.Add(1)
	//	go func() {
	//		for b := range blake2b_256_c {
	//			blake2b_256_d.Write(b)
	//		}
	//		wg.Done()
	//	}()
	//}
	//if hasHash(s_blake2b512) {
	//	wg.Add(1)
	//	go func() {
	//		for b := range blake2b_512_c {
	//			blake2b_512_d.Write(b)
	//		}
	//		wg.Done()
	//	}()
	//}
	//
	//data := make([]byte, 8192) // 8192 appears to be optimal
	//for {
	//	data = data[:cap(data)]
	//	n, err := file.Read(data)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//
	//		printError(fmt.Sprintf("reading file %s: %s", filename, err.Error()))
	//		return
	//	}
	//
	//	data = data[:n]
	//
	//	if hasHash(s_md5) {
	//		md5c <- data
	//	}
	//	if hasHash(s_sha1) {
	//		sha1c <- data
	//	}
	//	if hasHash(s_sha256) {
	//		sha256c <- data
	//	}
	//	if hasHash(s_sha512) {
	//		sha512c <- data
	//	}
	//	if hasHash(s_blake2b256) {
	//		blake2b_256_c <- data
	//	}
	//	if hasHash(s_blake2b512) {
	//		blake2b_512_c <- data
	//	}
	//}
	//
	//close(md5c)
	//close(sha1c)
	//close(sha256c)
	//close(sha512c)
	//close(blake2b_256_c)
	//close(blake2b_512_c)
	//
	//wg.Wait()
	//
	//fmt.Println(filename)
	//if hasHash(s_md5) {
	//	fmt.Println("        MD5 " + hex.EncodeToString(md5_d.Sum(nil)))
	//}
	//if hasHash(s_sha1) {
	//	fmt.Println("       SHA1 " + hex.EncodeToString(sha1_d.Sum(nil)))
	//}
	//if hasHash(s_sha256) {
	//	fmt.Println("     SHA256 " + hex.EncodeToString(sha256_d.Sum(nil)))
	//}
	//if hasHash(s_sha512) {
	//	fmt.Println("     SHA512 " + hex.EncodeToString(sha512_d.Sum(nil)))
	//}
	//if hasHash(s_blake2b256) {
	//	fmt.Println("Blake2b 256 " + hex.EncodeToString(blake2b_256_d.Sum(nil)))
	//}
	//if hasHash(s_blake2b512) {
	//	fmt.Println("Blake2b 512 " + hex.EncodeToString(blake2b_512_d.Sum(nil)))
	//}
	//fmt.Println("")
}


func processStandardInput(output chan Result) {
	total, nChunks := int64(0), int64(0)
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	md4_d := md4.New()
	md5_d := md5.New()
	sha1_d := sha1.New()
	sha256_d := sha256.New()
	sha512_d := sha512.New()
	blake2b_256_d := blake2b.New256()
	blake2b_512_d := blake2b.New512()
	sha3_224_d := sha3.New224()
	sha3_256_d := sha3.New256()
	sha3_384_d := sha3.New384()
	sha3_512_d := sha3.New512()

	md4c := make(chan []byte, 10)
	md5c := make(chan []byte, 10)
	sha1c := make(chan []byte, 10)
	sha256c := make(chan []byte, 10)
	sha512c := make(chan []byte, 10)
	blake2b_256_c := make(chan []byte, 10)
	blake2b_512_c := make(chan []byte, 10)
	sha3_224_c := make(chan []byte, 10)
	sha3_256_c := make(chan []byte, 10)
	sha3_384_c := make(chan []byte, 10)
	sha3_512_c := make(chan []byte, 10)

	var wg sync.WaitGroup

	if hasHash(HashNames.MD4) {
		wg.Add(1)
		go func() {
			for b := range md4c {
				md4_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.MD5) {
		wg.Add(1)
		go func() {
			for b := range md5c {
				md5_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA1) {
		wg.Add(1)
		go func() {
			for b := range sha1c {
				sha1_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA256) {
		wg.Add(1)
		go func() {
			for b := range sha256c {
				sha256_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA512) {
		wg.Add(1)
		go func() {
			for b := range sha512c {
				sha512_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Blake2b256) {
		wg.Add(1)
		go func() {
			for b := range blake2b_256_c {
				blake2b_256_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Blake2b512) {
		wg.Add(1)
		go func() {
			for b := range blake2b_512_c {
				blake2b_512_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Sha3224) {
		wg.Add(1)
		go func() {
			for b := range sha3_224_c {
				sha3_224_d.Write(b)
			}
			wg.Done()
		}()
	}
	if hasHash(HashNames.Sha3256) {
		wg.Add(1)
		go func() {
			for b := range sha3_256_c {
				sha3_256_d.Write(b)
			}
			wg.Done()
		}()
	}
	if hasHash(HashNames.Sha3384) {
		wg.Add(1)
		go func() {
			for b := range sha3_384_c {
				sha3_384_d.Write(b)
			}
			wg.Done()
		}()
	}
	if hasHash(HashNames.Sha3512) {
		wg.Add(1)
		go func() {
			for b := range sha3_512_c {
				sha3_512_d.Write(b)
			}
			wg.Done()
		}()
	}

	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if n == 0 {
			if err == nil {
				continue
			}

			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}

		nChunks++
		total += int64(len(buf))

		if hasHash(HashNames.MD4) {
			md4c <- buf
		}
		if hasHash(HashNames.MD5) {
			md5c <- buf
		}
		if hasHash(HashNames.SHA1) {
			sha1c <- buf
		}
		if hasHash(HashNames.SHA256) {
			sha256c <- buf
		}
		if hasHash(HashNames.SHA512) {
			sha512c <- buf
		}
		if hasHash(HashNames.Blake2b256) {
			blake2b_256_c <- buf
		}
		if hasHash(HashNames.Blake2b512) {
			blake2b_512_c <- buf
		}
		if hasHash(HashNames.Sha3224) {
			sha3_224_c <- buf
		}
		if hasHash(HashNames.Sha3256) {
			sha3_256_c <- buf
		}
		if hasHash(HashNames.Sha3384) {
			sha3_384_c <- buf
		}
		if hasHash(HashNames.Sha3512) {
			sha3_512_c <- buf
		}

		// process buf
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

	close(md4c)
	close(md5c)
	close(sha1c)
	close(sha256c)
	close(sha512c)
	close(blake2b_256_c)
	close(blake2b_512_c)
	close(sha3_224_c)
	close(sha3_256_c)
	close(sha3_384_c)
	close(sha3_512_c)

	wg.Wait()

	output <- Result{
		File:       "stdin",
		Bytes:      total,
		MD4:        hex.EncodeToString(md4_d.Sum(nil)),
		MD5:        hex.EncodeToString(md5_d.Sum(nil)),
		SHA1:       hex.EncodeToString(sha1_d.Sum(nil)),
		SHA256:     hex.EncodeToString(sha256_d.Sum(nil)),
		SHA512:     hex.EncodeToString(sha512_d.Sum(nil)),
		Blake2b256: hex.EncodeToString(blake2b_256_d.Sum(nil)),
		Blake2b512: hex.EncodeToString(blake2b_512_d.Sum(nil)),
		Sha3224:    hex.EncodeToString(sha3_224_d.Sum(nil)),
		Sha3256:    hex.EncodeToString(sha3_256_d.Sum(nil)),
		Sha3384:    hex.EncodeToString(sha3_384_d.Sum(nil)),
		Sha3512:    hex.EncodeToString(sha3_512_d.Sum(nil)),
	}

	close(output)
}

// For files over a certain size it is faster to process them using
// memory mapped files which this method does
// NB this does not play well with Windows as it will never
// be able to unmap the file "error unmapping: FlushFileBuffers: Access is denied."
func processMemoryMap(filename string) (Result, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		printError(fmt.Sprintf("opening file %s: %s", filename, err.Error()))
		return Result{}, err
	}

	mmap, err := mmapgo.Map(file, mmapgo.RDONLY, 0)
	if err != nil {
		printError(fmt.Sprintf("mapping file %s: %s", filename, err.Error()))
		return Result{}, err
	}

	md4_d := md4.New()
	md5_d := md5.New()
	sha1_d := sha1.New()
	sha256_d := sha256.New()
	sha512_d := sha512.New()
	blake2b_256_d := blake2b.New256()
	blake2b_512_d := blake2b.New512()
	sha3_224_d := sha3.New224()
	sha3_256_d := sha3.New256()
	sha3_384_d := sha3.New384()
	sha3_512_d := sha3.New512()

	md4c := make(chan []byte, 10)
	md5c := make(chan []byte, 10)
	sha1c := make(chan []byte, 10)
	sha256c := make(chan []byte, 10)
	sha512c := make(chan []byte, 10)
	blake2b_256_c := make(chan []byte, 10)
	blake2b_512_c := make(chan []byte, 10)
	sha3_224_c := make(chan []byte, 10)
	sha3_256_c := make(chan []byte, 10)
	sha3_384_c := make(chan []byte, 10)
	sha3_512_c := make(chan []byte, 10)

	var wg sync.WaitGroup

	if hasHash(HashNames.MD4) {
		wg.Add(1)
		go func() {
			for b := range md4c {
				md4_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.MD5) {
		wg.Add(1)
		go func() {
			for b := range md5c {
				md5_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA1) {
		wg.Add(1)
		go func() {
			for b := range sha1c {
				sha1_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA256) {
		wg.Add(1)
		go func() {
			for b := range sha256c {
				sha256_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA512) {
		wg.Add(1)
		go func() {
			for b := range sha512c {
				sha512_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Blake2b256) {
		wg.Add(1)
		go func() {
			for b := range blake2b_256_c {
				blake2b_256_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Blake2b512) {
		wg.Add(1)
		go func() {
			for b := range blake2b_512_c {
				blake2b_512_d.Write(b)
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Sha3224) {
		wg.Add(1)
		go func() {
			for b := range sha3_224_c {
				sha3_224_d.Write(b)
			}
			wg.Done()
		}()
	}
	if hasHash(HashNames.Sha3256) {
		wg.Add(1)
		go func() {
			for b := range sha3_256_c {
				sha3_256_d.Write(b)
			}
			wg.Done()
		}()
	}
	if hasHash(HashNames.Sha3384) {
		wg.Add(1)
		go func() {
			for b := range sha3_384_c {
				sha3_384_d.Write(b)
			}
			wg.Done()
		}()
	}
	if hasHash(HashNames.Sha3512) {
		wg.Add(1)
		go func() {
			for b := range sha3_512_c {
				sha3_512_d.Write(b)
			}
			wg.Done()
		}()
	}

	total := len(mmap)
	fileStartTime := makeTimestampMilli()

	// 1,048,576 = 2^20
	// No idea if this read size is optimal
	// TODO test out size here to find optimal for SSD
	for i := 0; i < total; i += 1048576 {
		end := i + 1048576
		if end > total {
			end = total
		}

		if hasHash(HashNames.MD4) {
			md4c <- mmap[i:end]
		}
		if hasHash(HashNames.MD5) {
			md5c <- mmap[i:end]
		}
		if hasHash(HashNames.SHA1) {
			sha1c <- mmap[i:end]
		}
		if hasHash(HashNames.SHA256) {
			sha256c <- mmap[i:end]
		}
		if hasHash(HashNames.SHA512) {
			sha512c <- mmap[i:end]
		}
		if hasHash(HashNames.Blake2b256) {
			blake2b_256_c <- mmap[i:end]
		}
		if hasHash(HashNames.Blake2b512) {
			blake2b_512_c <- mmap[i:end]
		}
		if hasHash(HashNames.Sha3224) {
			sha3_224_c <- mmap[i:end]
		}
		if hasHash(HashNames.Sha3256) {
			sha3_256_c <- mmap[i:end]
		}
		if hasHash(HashNames.Sha3384) {
			sha3_384_c <- mmap[i:end]
		}
		if hasHash(HashNames.Sha3512) {
			sha3_512_c <- mmap[i:end]
		}
	}

	if Trace {
		printTrace(fmt.Sprintf("milliseconds reading mmap file: %s: %d", filename, makeTimestampMilli()-fileStartTime))
	}

	close(md4c)
	close(md5c)
	close(sha1c)
	close(sha256c)
	close(sha512c)
	close(blake2b_256_c)
	close(blake2b_512_c)
	close(sha3_224_c)
	close(sha3_256_c)
	close(sha3_384_c)
	close(sha3_512_c)

	wg.Wait()

	if err := mmap.Unmap(); err != nil {
		printError(fmt.Sprintf("unmapping file %s: %s", filename, err.Error()))
	}

	return Result{
		File:       filename,
		Bytes:      int64(total),
		MD4:        hex.EncodeToString(md4_d.Sum(nil)),
		MD5:        hex.EncodeToString(md5_d.Sum(nil)),
		SHA1:       hex.EncodeToString(sha1_d.Sum(nil)),
		SHA256:     hex.EncodeToString(sha256_d.Sum(nil)),
		SHA512:     hex.EncodeToString(sha512_d.Sum(nil)),
		Blake2b256: hex.EncodeToString(blake2b_256_d.Sum(nil)),
		Blake2b512: hex.EncodeToString(blake2b_512_d.Sum(nil)),
		Sha3224:    hex.EncodeToString(sha3_224_d.Sum(nil)),
		Sha3256:    hex.EncodeToString(sha3_256_d.Sum(nil)),
		Sha3384:    hex.EncodeToString(sha3_384_d.Sum(nil)),
		Sha3512:    hex.EncodeToString(sha3_512_d.Sum(nil)),
	}, nil
}

// For files under a certain size its faster to just read them into memory in one
// chunk and then process them which this method does
// NB there is little point in multi-processing at this level, it would be
// better done on the input channel if required
func processReadFile(filename string) (Result, error) {
	startTime := makeTimestampNano()
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		printError(fmt.Sprintf("Unable to read file %s into memory with error %s", filename, err.Error()))
		return Result{}, err
	}

	if Trace {
		printTrace(fmt.Sprintf("nanoseconds reading file: %s: %d", filename, makeTimestampNano()-startTime))
	}

	var wg sync.WaitGroup
	result := Result{}

	if hasHash(HashNames.MD4) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := md4.New()
			d.Write(content)
			result.MD4 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing md4: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.MD5) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := md5.New()
			d.Write(content)
			result.MD5 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing md5: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA1) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha1.New()
			d.Write(content)
			result.SHA1 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha1: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA256) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha256.New()
			d.Write(content)
			result.SHA256 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha256: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.SHA512) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha512.New()
			d.Write(content)
			result.SHA512 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha512: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Blake2b256) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := blake2b.New256()
			d.Write(content)
			result.Blake2b256 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing blake2b-256: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Blake2b512) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := blake2b.New512()
			d.Write(content)
			result.Blake2b512 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing blake2b-512: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Sha3224) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha3.New224()
			d.Write(content)
			result.Sha3224 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha3-224: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Sha3256) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha3.New256()
			d.Write(content)
			result.Sha3256 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha3-256: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Sha3384) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha3.New384()
			d.Write(content)
			result.Sha3384 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha3-384: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	if hasHash(HashNames.Sha3512) {
		wg.Add(1)
		go func() {
			startTime = makeTimestampNano()
			d := sha3.New512()
			d.Write(content)
			result.Sha3512 = hex.EncodeToString(d.Sum(nil))

			if Trace {
				printTrace(fmt.Sprintf("nanoseconds processing sha3-512: %s: %d", filename, makeTimestampNano()-startTime))
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return result, nil
}
