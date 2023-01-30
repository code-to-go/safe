

add_library(safepool SHARED IMPORTED GLOBAL)
#set_property(TARGET safepool PROPERTY IMPORTED_NO_SONAME 1)
set_target_properties(safepool PROPERTIES IMPORTED_LOCATION ${CMAKE_CURRENT_SOURCE_DIR}/src/main/jniLibs/${ANDROID_ABI}/libsafepool.so)
#include_directories(${CMAKE_CURRENT_SOURCE_DIR}../../lib/mobile/jniLibs/${CMAKE_ANDROID_ARCH_ABI}/)
 