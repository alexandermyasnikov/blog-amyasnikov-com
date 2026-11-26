#include <stdint.h>
#include <stddef.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <sys/mman.h>

void print_hex_dump(const uint8_t *data, size_t size) {
    printf("Code dump (%zu bytes):\n", size);
    for (size_t i = 0; i < size; i++) {
        printf("%02x ", data[i]);
        if ((i + 1) % 16 == 0) printf("\n");
    }
    printf("\n");
}

int main() {
    const size_t code_size = 1024;

    void *exec_mem = mmap(NULL, code_size,
                          PROT_READ | PROT_WRITE | PROT_EXEC,
                          MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if (exec_mem == MAP_FAILED) {
        perror("mmap");
        return 1;
    }

    printf("Input shellcode (max %zu bytes):\n", code_size);
    ssize_t bytes_read = read(STDIN_FILENO, exec_mem, code_size);

    if (bytes_read <= 0) {
        perror("read failed");
        munmap(exec_mem, code_size);
        return 1;
    }

    print_hex_dump(exec_mem, bytes_read);

    int (*func)() = (int (*)())exec_mem;
    func();

    munmap(exec_mem, code_size);

    return 0;
}
