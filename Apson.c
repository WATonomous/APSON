import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
import requests
from bs4 import BeautifulSoup
import time

current_messages = []

def create_discussion(description):
    url = "https://api.github.com/WATonomous/infrastructure-support/discussions"  # Replace with actual repo details
    headers = {
        "Authorization": "Bearer YOUR_GITHUB_TOKEN",  # Replace with a valid GitHub token
        "Accept": "application/vnd.github.v3+json"
    }
    data = {
        "title": "Power Outage Notification - CPH Electrical Shutdown",
        "body": description,
        "category_id": ""  # Im not sure what this needs to be
    }
    
    try:
        response = requests.post(url, json=data, headers=headers)
        if response.status_code == 201:
            print("GitHub Discussion post created successfully!")
        else:
            print(f"Error creating GitHub Discussion post: {response.status_code}, {response.text}")
    except Exception as e:
        print(f"Exception occurred while creating GitHub Discussion post: {e}")
        
def send_email(des):
    email_subject = "Power Outage Notification - CPH Electrical Shutdown"
    email_body = "The following power outages have been detected:\n\n"
    email_body += des

    sender_email = "...@gmail.com" # pick an email
    sender_password = ""  # Use an app password 
    recipient_emails = ["...@gmail.com", "...@gmail.com"] # Change to the actual recipients

    msg = MIMEMultipart()
    msg["From"] = sender_email
    msg["To"] = ", ".join(recipient_emails)
    msg["Subject"] = email_subject
    msg.attach(MIMEText(email_body, "plain"))

    try:
        server = smtplib.SMTP("smtp.gmail.com", 587)  # Use correct SMTP server
        server.starttls()
        server.login(sender_email, sender_password)
        server.sendmail(sender_email, recipient_emails, msg.as_string())
        server.quit()
        print("Email sent successfully!")
    except Exception as e:
        print(f"Error sending email: {e}")

def get_info():
    url = "https://plantops.uwaterloo.ca/service-interruptions/" 
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"
    }

    response = requests.get(url, headers=headers)

    if response.status_code == 200:
        # Parse the HTML content
        soup = BeautifulSoup(response.content, "html.parser")
        
        # Find all <a> tags that contain the outage descriptions
        announcements = soup.find_all("a", class_="w3-leftbar")  # This assumes the outage descriptions are inside <a> tags with class="w3-leftbar"
        
        for announcement in announcements:
            # Print only the text within the <a> tag, which contains the description
            description = announcement.get_text(strip=True)  # strip=True removes extra whitespace
            if "CPH" in description and "electrical shutdown" in description.lower():
                if description in current_messages:
                    pass
                else:
                    print(description) # we probobly dont actually need this
                    send_email(description)          # we can switch this out for add discussion post
                    current_messages.append(description)
    else:
        print("er")

while True:
    get_info()
    time.sleep(86400)
