#include <stdio.h>
#include <unistd.h>
#include <sys/resource.h>
#include <sys/sysinfo.h>

struct CalcRAM {
  struct rusage usage;
  struct sysinfo sys;
  size_t start_memory;
  size_t end_memory;
  size_t total_memory;
  size_t memory_usage;
};

struct CalcRAM *new_CalcRAM() {
  struct CalcRAM *ram = malloc(sizeof(struct CalcRAM));
  getrusage(RUSAGE_SELF, &ram->usage);
  sysinfo(&ram->sys);
  ram->total_memory = ram->sys.totalram * ram->sys.mem_unit;
  return ram;
}

void start_measuring(struct CalcRAM *ram) {
  getrusage(RUSAGE_SELF, &ram->usage);
  ram->start_memory = ram->usage.ru_maxrss;
}

void stop_measuring(struct CalcRAM *ram) {
  getrusage(RUSAGE_SELF, &ram->usage);
  ram->end_memory = ram->usage.ru_maxrss;
  ram->memory_usage = ram->end_memory - ram->start_memory;
}

size_t calculate_memory_usage(struct CalcRAM *ram) {
  return ram->memory_usage;
}

float calculate_percent_memory_usage(struct CalcRAM *ram) {
  size_t memory_usage = calculate_memory_usage(ram);
  char percent_memory_usage_str[16];
  sprintf(percent_memory_usage_str, "%.10f", (float)memory_usage / (float)ram->total_memory * 100);
  float percent_memory_usage = atof(percent_memory_usage_str);
  return percent_memory_usage;
}

size_t get_total_memory(struct CalcRAM *ram) {
  return ram->total_memory;
}

size_t get_free_memory(struct CalcRAM *ram) {
  return ram->sys.freeram * ram->sys.mem_unit;
}

size_t get_used_memory(struct CalcRAM *ram) {
  return ram->total_memory - ram->sys.freeram * ram->sys.mem_unit;
}

float calculate_percent_free_memory(struct CalcRAM *ram) {
  size_t free_memory = get_free_memory(ram);
  return (float)free_memory / (float)ram->total_memory * 100;
}

float calculate_percent_used_memory(struct CalcRAM *ram) {
  size_t used_memory = get_used_memory(ram);
  return (float)used_memory / (float)ram->total_memory * 100;
}
void stop_calc(struct CalcRAM *ram) {
  free(ram);
}