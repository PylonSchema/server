Write-Output "\
    drop table channel_members;\
    drop table channels;\
    drop table invitation_channels;\
    drop table origins;\
    drop table refresh_tokens;\
    drop table socials;\
    drop table user_token_pairs;\
    drop table users;\" | mysql -u root -p pylon
