struct MemoryUsageData {
    atomic_size_t max_memory_usage;
    pthread_t thread;
    size_t total_memory;
    struct sysinfo sys;
};

void *measureMemoryUsage(void *arg) {
    // Cast the argument to a pointer to a MemoryUsageData struct
    struct MemoryUsageData *data = (struct MemoryUsageData *)arg;

    // Initialize the rusage struct
    struct rusage usage;

    // Loop indefinitely
    while (1) {
        // Get the current memory usage from getrusage
        getrusage(RUSAGE_SELF, &usage);
        // Update the maximum memory usage if necessary
        size_t current_memory_usage = usage.ru_maxrss;
        size_t old_max_memory_usage = atomic_load(&data->max_memory_usage);
        while (current_memory_usage > old_max_memory_usage &&
               !atomic_compare_exchange_weak(&data->max_memory_usage, &old_max_memory_usage, current_memory_usage)) {
            // If the update failed, try again with the updated value of old_max_memory_usage
        }

        // Update the total memory and sysinfo data
        data->total_memory = get_total_memory(&data->sys);
        get_sys_info(&data->sys);

        // Sleep for a short time to reduce CPU usage
        usleep(1000);  // Sleep for 1 millisecond
    }
    return NULL;
}
// Function to start the memory usage thread
int startMemoryUsageThread(struct MemoryUsageData *data) {
    // Initialize the MemoryUsageData struct
    atomic_store(&data->max_memory_usage, 0);

    // Create the continuous thread
    return pthread_create(&data->thread, NULL, measureMemoryUsage, data);
}

// Function to stop the memory usage thread
int stopMemoryUsageThread(struct MemoryUsageData *data) {
    // Stop the thread
    return pthread_cancel(data->thread);
}


size_t calculate_memory_usage(struct MemoryUsageData *data) {
  return data->max_memory_usage;
}

float calculate_percent_memory_usage(struct MemoryUsageData *data) {
  size_t memory_usage = calculate_memory_usage(data);
  char percent_memory_usage_str[16];
  sprintf(percent_memory_usage_str, "%.10f", (float)memory_usage / (float)data->total_memory * 100);
  float percent_memory_usage = atof(percent_memory_usage_str);
  return percent_memory_usage;
}

size_t get_total_memory(struct MemoryUsageData *data) {
  return data->total_memory;
}

size_t get_free_memory(struct MemoryUsageData *data) {
  return data->sys.freeram * data->sys.mem_unit;
}
size_t get_used_memory(struct MemoryUsageData *data) {
  return data->total_memory - data->sys.freeram * data->sys.mem_unit;
}
float calculate_percent_free_memory(struct MemoryUsageData *data) {
size_t free_memory = get_free_memory(data);
return (float)free_memory / (float)data->total_memory * 100;
}

float calculate_percent_used_memory(struct MemoryUsageData *data) {
size_t used_memory = get_used_memory(data);
return (float)used_memory / (float)data->total_memory * 100;
}

void stop_calc(struct MemoryUsageData *data) {
free(data);
}