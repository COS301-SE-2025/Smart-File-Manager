import 'package:flutter/material.dart';
import '../constants.dart';
import '../services/settings_service.dart';
import '../services/password_service.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  String _selectedNamingConvention = "CAMEL";
  String _selectedColorPreset = "Default";
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadSettings();
  }

  Future<void> _loadSettings() async {
    final settings = SettingsService.instance;
    final namingConvention = await settings.getNamingConvention();
    final colorPreset = await settings.getColorPreset();

    setState(() {
      _selectedNamingConvention = namingConvention;
      _selectedColorPreset = colorPreset;
      _isLoading = false;
    });
  }

  Future<void> _updateNamingConvention(String convention) async {
    await SettingsService.instance.setNamingConvention(convention);
    setState(() {
      _selectedNamingConvention = convention;
    });
  }

  Future<void> _updateColorPreset(String preset) async {
    await SettingsService.instance.setColorPreset(preset);
    setState(() {
      _selectedColorPreset = preset;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(
        child: CircularProgressIndicator(color: kprimaryColor),
      );
    }

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Settings',
              style: TextStyle(
                color: Colors.white,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            Divider(color: Color(0xff3D3D3D)),
            const SizedBox(height: 20),
            // Folder Naming Convention Section
            _buildSectionTitle('Folder Naming Convention'),
            const SizedBox(height: 16),
            _buildNamingConventionOptions(),

            const SizedBox(height: 40),

            // Color Presets Section
            _buildSectionTitle('Graph View Color Presets'),
            const SizedBox(height: 16),
            _buildColorPresetOptions(),

            const SizedBox(height: 40),

            // Security Section
            _buildSectionTitle('Security'),
            const SizedBox(height: 16),
            _buildSecurityOptions(),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Text(
      title,
      style: const TextStyle(
        color: Colors.white,
        fontSize: 18,
        fontWeight: FontWeight.w600,
      ),
    );
  }

  Widget _buildNamingConventionOptions() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xff242424),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: kOutlineBorder),
      ),
      child: Column(
        children:
            SettingsService.availableNamingConventions.map((convention) {
              return _buildCustomRadioTile(
                title: _formatConventionName(convention),
                value: convention,
                groupValue: _selectedNamingConvention,
                onChanged: (value) => _updateNamingConvention(value!),
              );
            }).toList(),
      ),
    );
  }

  Widget _buildColorPresetOptions() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xff242424),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: kOutlineBorder),
      ),
      child: Column(
        children:
            SettingsService.availablePresets.map((preset) {
              return _buildColorPresetTile(
                title: preset,
                value: preset,
                groupValue: _selectedColorPreset,
                onChanged: (value) => _updateColorPreset(value!),
              );
            }).toList(),
      ),
    );
  }

  Widget _buildCustomRadioTile({
    required String title,
    required String value,
    required String groupValue,
    required ValueChanged<String?> onChanged,
  }) {
    final isSelected = value == groupValue;

    return InkWell(
      onTap: () => onChanged(value),
      borderRadius: BorderRadius.circular(6),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
        child: Row(
          children: [
            Container(
              width: 20,
              height: 20,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(
                  color: isSelected ? kprimaryColor : const Color(0xff9CA3AF),
                  width: 2,
                ),
              ),
              child:
                  isSelected
                      ? Center(
                        child: Container(
                          width: 8,
                          height: 8,
                          decoration: const BoxDecoration(
                            shape: BoxShape.circle,
                            color: kprimaryColor,
                          ),
                        ),
                      )
                      : null,
            ),
            const SizedBox(width: 12),
            Text(
              title,
              style: TextStyle(
                color: isSelected ? Colors.white : const Color(0xff9CA3AF),
                fontSize: 14,
                fontWeight: isSelected ? FontWeight.w500 : FontWeight.normal,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildColorPresetTile({
    required String title,
    required String value,
    required String groupValue,
    required ValueChanged<String?> onChanged,
  }) {
    final isSelected = value == groupValue;
    final colors = _getPresetColors(value);

    return InkWell(
      onTap: () => onChanged(value),
      borderRadius: BorderRadius.circular(6),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
        child: Row(
          children: [
            Container(
              width: 20,
              height: 20,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(
                  color: isSelected ? kprimaryColor : const Color(0xff9CA3AF),
                  width: 2,
                ),
              ),
              child:
                  isSelected
                      ? Center(
                        child: Container(
                          width: 8,
                          height: 8,
                          decoration: const BoxDecoration(
                            shape: BoxShape.circle,
                            color: kprimaryColor,
                          ),
                        ),
                      )
                      : null,
            ),
            const SizedBox(width: 12),
            Text(
              title,
              style: TextStyle(
                color: isSelected ? Colors.white : const Color(0xff9CA3AF),
                fontSize: 14,
                fontWeight: isSelected ? FontWeight.w500 : FontWeight.normal,
              ),
            ),
            const Spacer(),
            // Color preview
            Row(
              children:
                  colors.map((color) {
                    return Container(
                      width: 16,
                      height: 16,
                      margin: const EdgeInsets.only(left: 4),
                      decoration: BoxDecoration(
                        color: color,
                        shape: BoxShape.circle,
                        border: Border.all(color: const Color(0xff3D3D3D)),
                      ),
                    );
                  }).toList(),
            ),
          ],
        ),
      ),
    );
  }

  String _formatConventionName(String convention) {
    switch (convention) {
      case "CAMEL":
        return "Camel Case (myFolderName)";
      case "SNAKE":
        return "Snake Case (my_folder_name)";
      case "PASCAL":
        return "Pascal Case (MyFolderName)";
      case "KEBAB":
        return "Kebab Case (my-folder-name)";
      case "SPACE":
        return "Space Case (My Folder Name)";
      default:
        return convention;
    }
  }

  List<Color> _getPresetColors(String preset) {
    return SettingsService.getPresetColors(preset).take(3).toList();
  }

  Widget _buildSecurityOptions() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xff242424),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: kOutlineBorder),
      ),
      child: Column(
        children: [
          _buildSecurityTile(
            title: 'Change Password',
            subtitle: 'Update your password',
            icon: Icons.lock_reset,
            onTap: _showChangePasswordDialog,
          ),
        ],
      ),
    );
  }

  Widget _buildSecurityTile({
    required String title,
    required String subtitle,
    required IconData icon,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(6),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(borderRadius: BorderRadius.circular(8)),
              child: Icon(icon, color: kprimaryColor, size: 20),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    subtitle,
                    style: const TextStyle(
                      color: Color(0xff9CA3AF),
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ),
            const Icon(Icons.chevron_right, color: Color(0xff9CA3AF), size: 20),
          ],
        ),
      ),
    );
  }

  void _showChangePasswordDialog() {
    showDialog(context: context, builder: (context) => _ChangePasswordDialog());
  }
}

class _ChangePasswordDialog extends StatefulWidget {
  @override
  State<_ChangePasswordDialog> createState() => _ChangePasswordDialogState();
}

class _ChangePasswordDialogState extends State<_ChangePasswordDialog> {
  final TextEditingController _currentPasswordController =
      TextEditingController();
  final TextEditingController _newPasswordController = TextEditingController();
  final TextEditingController _confirmPasswordController =
      TextEditingController();

  bool _isCurrentVisible = false;
  bool _isNewVisible = false;
  bool _isConfirmVisible = false;
  bool _isLoading = false;
  String _errorMessage = '';

  @override
  void dispose() {
    _currentPasswordController.dispose();
    _newPasswordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  Future<void> _changePassword() async {
    final currentPassword = _currentPasswordController.text;
    final newPassword = _newPasswordController.text;
    final confirmPassword = _confirmPasswordController.text;

    setState(() {
      _errorMessage = '';
      _isLoading = true;
    });

    // Validation
    if (currentPassword.isEmpty ||
        newPassword.isEmpty ||
        confirmPassword.isEmpty) {
      setState(() {
        _errorMessage = 'Please fill in all fields';
        _isLoading = false;
      });
      return;
    }

    if (newPassword.length < 4) {
      setState(() {
        _errorMessage = 'New password must be at least 4 characters';
        _isLoading = false;
      });
      return;
    }

    if (newPassword != confirmPassword) {
      setState(() {
        _errorMessage = 'New passwords do not match';
        _isLoading = false;
      });
      return;
    }

    // Change password
    final success = await PasswordService.instance.changePassword(
      currentPassword,
      newPassword,
    );

    setState(() {
      _isLoading = false;
    });

    if (success) {
      if (mounted) {
        Navigator.of(context).pop();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Password changed successfully'),
            backgroundColor: Colors.green,
          ),
        );
      }
    } else {
      setState(() {
        _errorMessage = 'Current password is incorrect';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      backgroundColor: const Color(0xff242424),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: const BorderSide(color: kOutlineBorder),
      ),
      child: Container(
        padding: const EdgeInsets.all(24),
        constraints: const BoxConstraints(maxWidth: 400),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.lock_reset, color: kprimaryColor, size: 24),
                const SizedBox(width: 12),
                const Text(
                  'Change Password',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Current password
            _buildPasswordField(
              label: 'Current Password',
              controller: _currentPasswordController,
              isVisible: _isCurrentVisible,
              onVisibilityToggle:
                  () => setState(() => _isCurrentVisible = !_isCurrentVisible),
              hintText: 'Enter current password',
            ),
            const SizedBox(height: 16),

            // New password
            _buildPasswordField(
              label: 'New Password',
              controller: _newPasswordController,
              isVisible: _isNewVisible,
              onVisibilityToggle:
                  () => setState(() => _isNewVisible = !_isNewVisible),
              hintText: 'Enter new password',
            ),
            const SizedBox(height: 16),

            // Confirm new password
            _buildPasswordField(
              label: 'Confirm New Password',
              controller: _confirmPasswordController,
              isVisible: _isConfirmVisible,
              onVisibilityToggle:
                  () => setState(() => _isConfirmVisible = !_isConfirmVisible),
              hintText: 'Confirm new password',
            ),

            // Error message
            if (_errorMessage.isNotEmpty) ...[
              const SizedBox(height: 12),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.red.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(6),
                  border: Border.all(color: Colors.red.withValues(alpha: 0.3)),
                ),
                child: Row(
                  children: [
                    Icon(
                      Icons.error_outline,
                      color: Colors.red.shade400,
                      size: 16,
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        _errorMessage,
                        style: TextStyle(
                          color: Colors.red.shade400,
                          fontSize: 12,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],

            const SizedBox(height: 24),

            // Buttons
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton(
                  onPressed:
                      _isLoading ? null : () => Navigator.of(context).pop(),
                  style: TextButton.styleFrom(
                    foregroundColor: const Color(0xff9CA3AF),
                  ),
                  child: const Text('Cancel'),
                ),
                const SizedBox(width: 12),
                ElevatedButton(
                  onPressed: _isLoading ? null : _changePassword,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: kprimaryColor,
                    foregroundColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(6),
                    ),
                  ),
                  child:
                      _isLoading
                          ? const SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.black,
                            ),
                          )
                          : const Text('Change Password'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPasswordField({
    required String label,
    required TextEditingController controller,
    required bool isVisible,
    required VoidCallback onVisibilityToggle,
    required String hintText,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 14,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            color: kScaffoldColor,
            borderRadius: BorderRadius.circular(6),
            border: Border.all(color: kOutlineBorder),
          ),
          child: TextField(
            controller: controller,
            obscureText: !isVisible,
            style: const TextStyle(color: Colors.white),
            decoration: InputDecoration(
              hintText: hintText,
              hintStyle: const TextStyle(color: Color(0xff9CA3AF)),
              border: InputBorder.none,
              contentPadding: const EdgeInsets.symmetric(
                horizontal: 12,
                vertical: 12,
              ),
              suffixIcon: IconButton(
                icon: Icon(
                  isVisible ? Icons.visibility : Icons.visibility_off,
                  color: const Color(0xff9CA3AF),
                  size: 20,
                ),
                onPressed: onVisibilityToggle,
              ),
            ),
          ),
        ),
      ],
    );
  }
}
