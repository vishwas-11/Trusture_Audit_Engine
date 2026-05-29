// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Trusture {
    event DonationReceived(
        address indexed donor,
        address indexed ngo,
        uint256 amount,
        string purpose,
        bytes32 proofHash,
        uint256 timestamp
    );

    function donate(
        address ngo,
        string calldata purpose,
        bytes32 proofHash
    ) external payable {
        require(msg.value > 0, "Donation must be > 0");
        require(ngo != address(0), "Invalid NGO");

        emit DonationReceived(
            msg.sender,
            ngo,
            msg.value,
            purpose,
            proofHash,
            block.timestamp
        );
    }
}
