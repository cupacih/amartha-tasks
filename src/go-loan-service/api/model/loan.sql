-- CREATE DATABASE loanservice;

CREATE TABLE Loan (
    ID VARCHAR(255) PRIMARY KEY,
    LoanID VARCHAR(255) NOT NULL,
    BorrowerID VARCHAR(255) NOT NULL,
    PrincipalAmount DECIMAL(18,2) NOT NULL,
    Rate DECIMAL(18,2) NOT NULL,
    ROI DECIMAL(18,2) NOT NULL,
    AgreementLink VARCHAR(255),
    State ENUM('proposed', 'approved', 'invested', 'disbursed') NOT NULL,
    InvestedAmount DECIMAL(18,2) NOT NULL DEFAULT 0,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_loan_id (LoanID),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;