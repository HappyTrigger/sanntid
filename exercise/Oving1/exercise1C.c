#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>


int var = 0;

void *thread_1(){
	int i = 0;
	for(i;i<1000;i++){
		var++;
	}
}

void *thread_2(){
	int i=0;
	printf("%s\n","hehehehe" );
	for(i; i<1000;i++){
		var--;
	}
}

int main(void){
	pthread_t thread_1;
	pthread_t thread_2;
	printf("%d",var);
	return 0;
}