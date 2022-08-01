#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <arpa/inet.h>

#define DATA_SIZE 1024

extern int bench_main(int argc, char *argv[]);

static const char *HTTP_header = "HTTP/1.1 %d %s\r\n"
      "Content-Type: text/plain\r\n"
			"Connection: " "Keep-Alive" "\r\n"
			"Date: Mon, 01 Jan 1970 00:00:00 GMT GMT\r\n"
			"Content-Length: %ld\r\n"
			"\r\n";
  
static void create_socket(int *sockfd, int *connfd)
{
	int port = 8070;
	struct sockaddr_in server_addr, cli;
	socklen_t addr_size;

  *sockfd = socket(AF_INET, SOCK_STREAM, 0);
	if(*sockfd < 0) {
		perror("Error in socket");
		exit(*sockfd);
	}
  printf("Socket created\n");

	bzero(&server_addr, sizeof(server_addr));

	server_addr.sin_family = AF_INET;
	server_addr.sin_port = htons(port);
	server_addr.sin_addr.s_addr = htonl(INADDR_ANY);

	if ((bind(*sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr))) != 0) {
		perror("socket bind failed...\n");
		exit(1);
	}
  printf("Address bound to socket\n");

	if (listen(*sockfd, 5) != 0) {
		perror("Listen failed...\n");
		exit(1);
	}

	printf("Waiting for connections\n");
	*connfd = accept(*sockfd, (struct sockaddr *)&cli, &addr_size);
	if (*connfd < 0) {
		perror("server accept failed...\n");
		exit(1);
	}

  printf("Connected to client\n");
}

static void send_results(int *sockfd, int *connfd)
{
	FILE *file;
  long output_size;
	char data[DATA_SIZE] = {0};
  char response[4*DATA_SIZE] = {0};

  printf("Sending results to one client\n");

	file = fopen("output.txt", "r");
	if (file == NULL) {
		exit(1);
	}

  // Get the size of the file
  fseek(file, 0L, SEEK_END);
  output_size = ftell(file);
  fseek(file, 0L, SEEK_SET);

  // Append HTTP/1.1 header
  snprintf(response, sizeof(response), HTTP_header, 200, "OK", output_size);

  // Also copy the whole file into the response
	while (fgets(data, DATA_SIZE, file)) {
		printf("%s", data);
    strncat(response, data, DATA_SIZE);
		bzero(data, DATA_SIZE);
	}

  // Send the response
  if (send(*connfd, response, strnlen(response, DATA_SIZE), 0) == -1) {
    perror("sending data failed...\n");
    exit(1);
  }

	close(*sockfd);
	fclose(file);
}

int main(int argc, char *argv[])
{
	int ret;
	int sockfd, connfd;

  create_socket(&sockfd, &connfd);

	ret = bench_main(argc, argv);
	if (ret != 0) {
		return ret;
	}

	send_results(&sockfd, &connfd);

	return 0;
}
