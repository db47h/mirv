#include <unistd.h>

int main(void) {
	write(1, "Hello, World!", 13);
	return 0;
}
