#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>


typedef struct item
{
	int amount;
	char* serial;
	char* dateExported;
	char* dateImported;
	char* expirationDate;
}Item;

/*implement LRU cache*/

//nodes for doubly linked list
typedef struct Node {
    struct Node *prev, *next;
    unsigned pageNumber; // the page number stored in this QNode
}Node;


// A Queue (A FIFO collection of Queue Nodes)
typedef struct Queue {
    unsigned count; // Number of filled frames
    unsigned numberOfFrames; // total number of frames
    Node *front, *rear;
}Queue;

typedef struct Hash
{
	int capacity;
	Node** array; 
}Hash;

Node* newNode(unsigned pageNumber)
{
	Node* temp = (Node*)malloc(sizeof(Node));
	temp->pageNumber = pageNumber;

	temp->prev = temp->next = NULL;

	return temp;
}

Queue* createQueue(int numberOfFrames)
{
	Queue* queue = (Queue*)malloc(sizeof(Queue));

	queue->count = 0;
	queue->numberOfFrames = numberOfFrames;
	queue->front = queue->rear = NULL;

	return queue;
}


Hash* createHash(int capacity)
{
	Hash* hash = (Hash*)malloc(sizeof(Hash));
	hash->capacity = capacity;

	hash->array = (Node**)malloc(hash->capacity * sizeof(Node**));

	int i;
	for(i = 0;i<hash->capacity;i++)
	{
		hash->array[i] = NULL;
	}
	return hash;
}

int IsQueueFull(Queue* queue)
{
	return queue->count == queue->numberOfFrames;
}

int IsQueueEmpty(Queue* queue)
{
	return queue->rear == NULL;
}

void dequeue(Queue* queue)
{
	if(IsQueueEmpty(queue))
	{
		return;
	}
	if (queue->front == queue->rear)
	{
		queue->front = NULL;
	}
	Node* temp = queue->rear;
	queue->rear->prev = queue->rear;
	if (queue->rear)
	{
		queue->rear->next = NULL;
	}
	free(temp);
	queue->count-=1;

}

void Enqueue(Queue* queue,Hash* hash, unsigned pageNumber)
{
	if(IsQueueFull(queue))
	{
		hash->array[queue->rear->pageNumber] = NULL;
        dequeue(queue);
	}

	Node* temp = newNode(pageNumber);
	temp->next = queue->front;

	if (IsQueueEmpty)
	{
		queue->front = queue->rear = temp;
	}
	else
	{
		queue->front->prev = temp;
		queue->front = temp;
	}
	hash->array[pageNumber] = temp;
	queue->count += 1;
}



