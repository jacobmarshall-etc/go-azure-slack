FROM ubuntu:trusty
ADD go-azure-slack .
CMD ["./go-azure-slack"]