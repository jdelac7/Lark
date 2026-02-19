plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

val natives by configurations.creating

android {
    namespace = "com.lark.app"
    compileSdk = 34

    defaultConfig {
        applicationId = "com.lark.app"
        minSdk = 26
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"

        // Server URL - override via gradle property or BuildConfig
        buildConfigField("String", "SERVER_URL", "\"http://10.0.2.2:9292\"")
    }

    buildTypes {
        release {
            isMinifyEnabled = false
        }
    }

    buildFeatures {
        compose = true
        buildConfig = true
    }

    composeOptions {
        kotlinCompilerExtensionVersion = "1.5.8"
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = "17"
        freeCompilerArgs += "-opt-in=androidx.compose.material3.ExperimentalMaterial3Api"
    }
}

dependencies {
    // Compose BOM
    val composeBom = platform("androidx.compose:compose-bom:2024.02.00")
    implementation(composeBom)

    // Compose
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    implementation("androidx.activity:activity-compose:1.8.2")
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.7.0")
    implementation("androidx.navigation:navigation-compose:2.7.6")

    // libGDX
    val gdxVersion = "1.13.0"
    implementation("com.badlogicgames.gdx:gdx:$gdxVersion")
    implementation("com.badlogicgames.gdx:gdx-backend-android:$gdxVersion")
    natives("com.badlogicgames.gdx:gdx-platform:$gdxVersion:natives-armeabi-v7a")
    natives("com.badlogicgames.gdx:gdx-platform:$gdxVersion:natives-arm64-v8a")
    natives("com.badlogicgames.gdx:gdx-platform:$gdxVersion:natives-x86")
    natives("com.badlogicgames.gdx:gdx-platform:$gdxVersion:natives-x86_64")

    // Retrofit
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.retrofit2:converter-gson:2.9.0")
    implementation("com.squareup.okhttp3:logging-interceptor:4.12.0")

    // Coroutines
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3")

    // Core
    implementation("androidx.core:core-ktx:1.12.0")
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.7.0")

    // Debug
    debugImplementation("androidx.compose.ui:ui-tooling")
}

// Extract libGDX native .so files into jniLibs so they get packaged into the APK
tasks.register("copyNatives") {
    doFirst {
        val jniLibsDir = file("src/main/jniLibs")
        jniLibsDir.deleteRecursively()
        natives.files.forEach { jar ->
            val abi = when {
                jar.name.contains("arm64-v8a") -> "arm64-v8a"
                jar.name.contains("armeabi-v7a") -> "armeabi-v7a"
                jar.name.contains("x86_64") -> "x86_64"
                jar.name.contains("x86") -> "x86"
                else -> return@forEach
            }
            val targetDir = File(jniLibsDir, abi)
            targetDir.mkdirs()
            zipTree(jar).matching { include("*.so") }.forEach { soFile ->
                soFile.copyTo(File(targetDir, soFile.name), overwrite = true)
            }
        }
    }
}

tasks.named("preBuild") {
    dependsOn("copyNatives")
}
