section .text

fun2:
    mov rdi, 0x0068732f6e69622f   ; '/bin/sh'
    push rdi

    mov rax, 59        ; sys_execve
    mov rdi, rsp       ; filename
    xor rsi, rsi       ; argv
    xor rdx, rdx       ; envp
    syscall
