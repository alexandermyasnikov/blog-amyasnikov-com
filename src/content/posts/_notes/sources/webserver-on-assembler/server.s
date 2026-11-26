.intel_syntax noprefix
.globl _start

.equ SYS_read, 0
.equ SYS_write, 1
.equ SYS_open, 2
.equ SYS_close, 3
.equ SYS_socket, 41
.equ SYS_accept, 43
.equ SYS_bind, 49
.equ SYS_listen, 50
.equ SYS_fork, 57
.equ SYS_exit, 60

.equ EXIT_OK, 0
.equ EXIT_ERROR, 1
.equ AF_INET, 2
.equ SOCK_STREAM, 1
.equ IPPROTO_IP, 0
.equ INADDR_ANY, 0
.equ O_RDONLY, 0
.equ O_WRONLY, 1
.equ O_CREAT, 64

.equ PORT, 8080


.section .data

response_header_ok:
    .ascii "HTTP/1.0 200 OK\r\n\r\n"
response_header_ok_length = . - response_header_ok


.section .text

_start:
    # vars:
    #   [rbp-16]   - 16b, sockaddr
    #   [rbp-24]   - 8b, srv_socket_fd
    #   [rbp-32]   - 8b, cli_socket_fd

    push rbp
    mov rbp, rsp
    sub rsp, 32


    # int socket(int domain, int type, int protocol);
    mov rdi, AF_INET       # domain
    mov rsi, SOCK_STREAM   # type
    mov rdx, IPPROTO_IP    # protocol
    mov rax, SYS_socket    # socket
    syscall

    test rax, rax
    js _exit_with_error


    mov [rbp-24], rax   # save srv_socket_fd


    # init sockaddr_in
    movw [rbp-16],   AF_INET
    movw ax, PORT
    xchg al, ah
    movw [rbp-16+2], ax
    movd [rbp-16+4], INADDR_ANY


    # int bind(int socket, const struct sockaddr *address, socklen_t address_len);
    mov rdi, [rbp-24]      # socket
    lea rsi, [rbp-16]      # address
    mov rdx, 16            # address_len
    mov rax, SYS_bind      # socket
    syscall

    test rax, rax
    js _exit_with_error


    # int listen(int socket, int backlog);
    mov rdi, [rbp-24]      # socket
    mov rsi, 0             # backlog
    mov rax, SYS_listen    # listen
    syscall

    test rax, rax
    js _exit_with_error


.accept:
    # int accept(int socket, struct sockaddr *restrict address, socklen_t *restrict address_len);
    mov rdi, [rbp-24]        # socket
    mov rsi, 0               # address
    mov rdx, 0               # adress_len
    mov rax, SYS_accept      # accept
    syscall

    test rax, rax
    js _exit_with_error


    mov [rbp-32], rax   # save cli_socket_fd


    # pid_t fork(void);
    mov rax, SYS_fork   # fork
    syscall


    test rax, rax
    js _exit_with_error


    mov r14, rax   # pid

    # close cli_socket_fd in parent
    # close srv_socket_fd in child
    jz .accept_set_srv_socket_fd
    mov r15, [rbp-32]
    jmp .accept_close_fd
.accept_set_srv_socket_fd:
    mov r15, [rbp-24]
.accept_close_fd:
    # int close(int fd);
    mov rdi, r15         # fd
    mov rax, SYS_close   # close
    syscall


    # parent process wait new clients
    test r14, r14
    jg .accept


    # child process precesses client
    mov rdi, [rbp-32]
    call process_request


    # exit(EXIT_OK)
    mov rdi, EXIT_OK
    mov rax, SYS_exit
    syscall


    mov rsp, rbp
    pop rbp
    ret


_exit_with_error:
    # exit(EXIT_ERROR)
    mov rdi, EXIT_ERROR
    mov rax, SYS_exit
    syscall


process_request:
    # args:
    #   rdi = fd
    # vars:
    #   [rbp-8]    - 8b, fd
    #   [rbp-16]   - 8b, buffer_length
    #   [rbp-1040] - 1024b, buffer

    push rbp
    mov rbp, rsp
    sub rsp, 1040


    mov [rbp-8], rdi   # save client socket fd


    # ssize_t read(int fd, void buf[.count], size_t count);
    mov rdi, [rbp-8]      # fd
    lea rsi, [rbp-1040]   # buf
    mov rdx, 1024         # count
    mov rax, SYS_read     # read
    syscall


    cmpb [rbp-1040], 'P'
    je _process_post
    mov rdi, [rbp-8]
    lea rsi, [rbp-1040]
    call process_get
    jmp _endif
_process_post:
    mov rdi, [rbp-8]
    lea rsi, [rbp-1040]
    mov rdx, rax
    call process_post
_endif:


    mov rsp, rbp
    pop rbp
    ret


process_post:
    # args:
    #   rdi = socket_fd
    #   rsi = &buffer
    #   rdx = buffer_length
    # vars:
    #   [rbp-8]    - 8b, socket_fd
    #   [rbp-16]   - 8b, &req_buffer
    #   [rbp-24]   - 8b, buffer_length
    #   [rbp-32]   - 8b, file_fd

    push rbp
    mov rbp, rsp
    sub rsp, 32

    mov [rbp-8], rdi
    mov [rbp-16], rsi
    mov [rbp-24], rdx


    # set \x00 after "xxx" in "POST /xxx HTTP/1.1\r\n....."
    mov r8, rsi
    add r8, 5
.process_post_path:
    add r8, 1
    cmpb [r8], ' '
    jne .process_post_path
.process_post_path_end:
    movb byte [r8-1], 0


    # int open(const char *pathname, int flags, /* mode_t mode */);
    lea rdi, [rsi]+6    # pathname
    mov rsi, O_WRONLY   # flags
    xor rsi, O_CREAT    # flags
    mov rdx, 0777       # mode
    mov rax, SYS_open   # open
    syscall

    test rax, rax
    js _exit_with_error


    mov [rbp-32], rax   # save fd


    # parse data between '\n' and end of body
    mov r8, [rbp-16]
    add r8, [rbp-24]
.process_body:
    sub r8, 1
    cmpb [r8], '\n'
    jne .process_body
.process_body_end:
    add r8, 1   # data_offset

    mov r9, [rbp-16]
    add r9, [rbp-24]
    sub r9, r8   # data_length


    # ssize_t write(int fd, void buf[.count], size_t count);
    mov rdi, [rbp-32]    # fd
    mov rsi, r8          # buf
    mov rdx, r9          # count
    mov rax, SYS_write   # write
    syscall


    # int close(int fd);
    mov rdi, [rbp-32]   # fd
    mov rax, SYS_close  # close
    syscall

    test rax, rax
    js _exit_with_error


    # ssize_t write(int fd, const void buf[.count], size_t count);
    mov rdi, [rbp-8]                     # fd
    lea rsi, [response_header_ok]        # buf
    mov rdx, response_header_ok_length   # count
    mov rax, SYS_write                   # write
    syscall


    mov rsp, rbp
    pop rbp
    ret

process_get:
    # args:
    #   rdi = socket_fd
    #   rsi = &buffer
    # vars:
    #   [rbp-8]    - 8b, socket_fd
    #   [rbp-16]   - 8b, &req_buffer
    #   [rbp-24]   - 8b, file_fd
    #   [rbp-1048] - 1024b, file_buffer
    #   [rbp-1056] - 8b, file_buffer_length

    push rbp
    mov rbp, rsp
    sub rsp, 1056

    mov [rbp-8], rdi
    mov [rbp-16], rsi


    # set \x00 after "xxx" in "GET /xxx HTTP/1.1\r\n....."
    mov r8, rsi
    add r8, 4
.process_get_path:
    add r8, 1
    cmpb [r8], ' '
    jne .process_get_path
.process_get_path_end:
    movb byte [r8-1], 0


    # int open(const char *pathname, int flags, /* mode_t mode */);
    lea rdi, [rsi]+5    # pathname
    mov rsi, O_RDONLY   # flags
    mov rdx, 0          # mode
    mov rax, SYS_open   # open
    syscall

    test rax, rax
    js _exit_with_error

    mov [rbp-24], rax


    # ssize_t read(int fd, void buf[.count], size_t count);
    mov rdi, [rbp-24]     # fd
    lea rsi, [rbp-1048]   # buf
    mov rdx, 1024         # count
    mov rax, SYS_read     # read
    syscall

    mov [rbp-1056], rax


    # int close(int fd);
    mov rdi, [rbp-24]   # fd
    mov rax, SYS_close  # close
    syscall


    # ssize_t write(int fd, const void buf[.count], size_t count);
    mov rdi, [rbp-8]                     # fd
    lea rsi, [response_header_ok]        # buf
    mov rdx, response_header_ok_length   # count
    mov rax, SYS_write                   # write
    syscall


    # ssize_t write(int fd, const void buf[.count], size_t count);
    mov rdi, [rbp-8]      # fd
    lea rsi, [rbp-1048]   # buf
    mov rdx, [rbp-1056]   # count
    mov rax, SYS_write    # write
    syscall
     

    mov rsp, rbp
    pop rbp
    ret
