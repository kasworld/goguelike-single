# prime number finder 

소수를 여러가지 방법(6가지)으로 구하는 패키지 입니다. 

# 기능 

MakePrimes 

    가장 기본적인 구한 소수들을 테이블에 보관해가면서 더 큰 소수들을 찾아 냅니다. 
    아래의 PrimeIntList 를 쓰지 않습니다. 


type PrimeIntList 

    구한 소수 테이블을 유지하면서 추가로 소수들을 구해가기 위한 정보들 담고 있습니다. 
    AppendFindTo : 주어진 인자 까지 수수를 찾아 냅니다. 
    MultiAppendFindTo : go channel 과 go routine worker 를 사용해서 multithread로 소수를 구합니다. 
    MultiAppendFindTo2 : channel 사용를 줄여 속도를 올린 함수 입니다. 
    MultiAppendFindTo3 : channel사용을 없애고 worker 별 계산 결과를 merge sort로 취합합니다. 
    MultiAppendFindTo4 : 계산용 buffer를 미리 alloc해서 slice의 크기를 늘이기 위한 오버헤드를 줄였습니다. 

# 사용 예제 

example/example.go 참조.