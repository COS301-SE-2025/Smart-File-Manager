import 'package:flutter/material.dart';
import '../services/password_service.dart';
import '../custom_widgets/password_setup_dialog.dart';
import '../custom_widgets/password_entry_dialog.dart';
import '../navigation/shell.dart';
import '../constants.dart';
import 'dart:io';

class AuthWrapper extends StatefulWidget {
  const AuthWrapper({super.key});

  @override
  State<AuthWrapper> createState() => _AuthWrapperState();
}

class _AuthWrapperState extends State<AuthWrapper> {
  bool _isLoading = true;
  bool _isAuthenticated = false;

  @override
  void initState() {
    super.initState();
    _checkAuthenticationStatus();
  }

  Future<void> _checkAuthenticationStatus() async {
    try {
      final isPasswordSet = await PasswordService.instance.isPasswordSet();
      
      if (!isPasswordSet) {
        // First time setup - show password setup dialog
        await _showPasswordSetupDialog();
      } else {
        // Password exists - show password entry dialog
        await _showPasswordEntryDialog();
      }
    } catch (e) {
      // Error occurred - show setup dialog as fallback
      await _showPasswordSetupDialog();
    }
  }

  Future<void> _showPasswordSetupDialog() async {
    while (!_isAuthenticated && mounted) {
      final result = await showDialog<bool>(
        context: context,
        barrierDismissible: false,
        builder: (context) => const PasswordSetupDialog(),
      );

      if (result == true) {
        setState(() {
          _isAuthenticated = true;
          _isLoading = false;
        });
      } else {
        // User cancelled setup - exit app
        _exitApp();
        return;
      }
    }
  }

  Future<void> _showPasswordEntryDialog() async {
    while (!_isAuthenticated && mounted) {
      final result = await showDialog<bool>(
        context: context,
        barrierDismissible: false,
        builder: (context) => const PasswordEntryDialog(),
      );

      if (result == true) {
        setState(() {
          _isAuthenticated = true;
          _isLoading = false;
        });
      } else {
        // User cancelled or chose to exit - exit app
        _exitApp();
        return;
      }
    }
  }

  void _exitApp() {
    exit(0);
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        backgroundColor: kScaffoldColor,
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 80,
                height: 80,
                decoration: BoxDecoration(
                  color: kprimaryColor.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: const Icon(
                  Icons.security,
                  size: 40,
                  color: kprimaryColor,
                ),
              ),
              const SizedBox(height: 24),
              const Text(
                'Smart File Manager',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              const Text(
                'Initializing security...',
                style: TextStyle(
                  color: Color(0xff9CA3AF),
                  fontSize: 14,
                ),
              ),
              const SizedBox(height: 32),
              const CircularProgressIndicator(
                color: kprimaryColor,
                strokeWidth: 2,
              ),
            ],
          ),
        ),
      );
    }

    if (_isAuthenticated) {
      return const Shell();
    }

    // Fallback - should not reach here
    return Scaffold(
      backgroundColor: kScaffoldColor,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.error_outline,
              size: 64,
              color: Colors.red,
            ),
            const SizedBox(height: 16),
            const Text(
              'Authentication Error',
              style: TextStyle(
                color: Colors.white,
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              'Unable to authenticate. Please restart the application.',
              style: TextStyle(
                color: Color(0xff9CA3AF),
                fontSize: 14,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: _exitApp,
              style: ElevatedButton.styleFrom(
                backgroundColor: kprimaryColor,
                foregroundColor: Colors.black,
              ),
              child: const Text('Exit Application'),
            ),
          ],
        ),
      ),
    );
  }
}