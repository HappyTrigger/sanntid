#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>


pthread_mutex_t mutex;
int var = 1;

void *thread_1(){
	int i = 0;
	for(i;i<10000000;i++){
		pthread_mutex_lock(&mutex);
		var++;
		pthread_mutex_unlock(&mutex);
	}
}

void *thread_2(){
	int i=0;
	for(i; i<10000000;i++){
		pthread_mutex_lock(&mutex);
		var--;
		pthread_mutex_unlock(&mutex);
	}
}

int main(void){
	pthread_t thread_12;
	pthread_t thread_21;

	pthread_create(&thread_12, NULL, thread_1, NULL);
	pthread_create(&thread_21, NULL, thread_2, NULL);
	
	pthread_join(thread_21, NULL);
	pthread_join(thread_12,NULL);
	printf("%d\n",var);
	return 0;
}