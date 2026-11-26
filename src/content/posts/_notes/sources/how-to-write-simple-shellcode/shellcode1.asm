section .text

fun1:
    push rbp
    mov rbp, rsp
    pop rbp

    mov rax, 0x00000a6f6c6c6568   ; 'hello'
    push rax

    mov rax, 1          ; sys_write
    mov rdi, 1          ; stdout
    mov rsi, rsp        ; buf
    mov rdx, 6          ; count
    syscall

    mov rax, 0

    mov rsp, rbp
    pop rbp
    ret
